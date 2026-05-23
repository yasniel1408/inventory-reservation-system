package reservations

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"inventory-reservation-system/backend/internal/store"
)

type PostgresStore struct {
	db *gorm.DB
}

func NewPostgresStore(db *gorm.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) ExpireActiveReservations(ctx context.Context) error {
	return store.WithTransaction(ctx, s.db, func(tx *gorm.DB) error {
		return expireActiveReservations(tx)
	})
}

func (s *PostgresStore) CreateReservation(ctx context.Context, input CreateInput, requestHash string) (CreateResult, error) {
	var result CreateResult
	var domainErr error

	err := store.WithTransaction(ctx, s.db, func(tx *gorm.DB) error {
		if err := expireActiveReservations(tx); err != nil {
			return err
		}

		scope := "POST:/reservations:" + input.SessionID
		insert := tx.Exec(`
			INSERT INTO idempotency_keys (scope, key, request_hash, status)
			VALUES (?, ?, ?, 'in_progress')
			ON CONFLICT (scope, key) DO NOTHING
		`, scope, input.IdempotencyKey, requestHash)
		if insert.Error != nil {
			return insert.Error
		}

		if insert.RowsAffected == 0 {
			replay, err := s.replayExistingIdempotency(tx, scope, input.IdempotencyKey, requestHash)
			if err != nil {
				return err
			}
			result = replay
			return nil
		}

		reservation, err := s.createFreshReservation(tx, input)
		if err != nil {
			code, apiCode, deterministic := idempotentErrorPayload(err)
			if deterministic {
				if updateErr := markIdempotencyFailed(tx, scope, input.IdempotencyKey, code, apiCode); updateErr != nil {
					return updateErr
				}
				domainErr = err
				return nil
			}
			return err
		}

		if err := markIdempotencyCompleted(tx, scope, input.IdempotencyKey, reservation); err != nil {
			return err
		}

		result = CreateResult{
			StatusCode:  http.StatusCreated,
			Reservation: reservation,
		}
		return nil
	})
	if domainErr != nil {
		return result, domainErr
	}

	return result, err
}

func (s *PostgresStore) ListActiveReservations(ctx context.Context, sessionID string) ([]Reservation, error) {
	var result []Reservation
	err := s.db.WithContext(ctx).Raw(`
		SELECT
			r.id::text AS id,
			r.item_id::text AS item_id,
			i.name AS item_name,
			r.session_id,
			r.quantity,
			r.status,
			r.expires_at,
			r.created_at,
			r.released_at,
			r.expired_at
		FROM reservations r
		INNER JOIN items i ON i.id = r.item_id
		WHERE r.session_id = ?
		  AND r.status = 'active'
		  AND r.expires_at > NOW()
		ORDER BY r.created_at DESC
	`, sessionID).Scan(&result).Error
	return result, err
}

func (s *PostgresStore) ReleaseReservation(ctx context.Context, sessionID string, reservationID string) (ReleaseResult, error) {
	var result ReleaseResult

	err := store.WithTransaction(ctx, s.db, func(tx *gorm.DB) error {
		if err := expireActiveReservations(tx); err != nil {
			return err
		}

		reservation, err := lockReservation(tx, sessionID, reservationID)
		if err != nil {
			return err
		}

		if reservation.Status != StatusActive {
			result = ReleaseResult{
				Reservation:   reservation,
				StockReturned: false,
				Noop:          true,
			}
			return nil
		}

		var released struct {
			ReleasedAt time.Time
		}
		if err := tx.Raw(`
			UPDATE reservations
			SET status = 'released',
			    released_at = NOW(),
			    updated_at = NOW()
			WHERE id = ?::uuid
			  AND status = 'active'
			RETURNING released_at
		`, reservation.ID).Scan(&released).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			UPDATE items
			SET reserved_stock = reserved_stock - ?,
			    updated_at = NOW()
			WHERE id = ?::uuid
		`, reservation.Quantity, reservation.ItemID).Error; err != nil {
			return err
		}

		reservation.Status = StatusReleased
		reservation.ReleasedAt = &released.ReleasedAt
		result = ReleaseResult{
			Reservation:   reservation,
			StockReturned: true,
			Noop:          false,
		}
		return nil
	})

	return result, err
}

func (s *PostgresStore) replayExistingIdempotency(tx *gorm.DB, scope string, key string, requestHash string) (CreateResult, error) {
	var row struct {
		RequestHash   string
		Status        string
		ResponseCode  sql.NullInt64
		ResponseBody  []byte
		ReservationID sql.NullString
	}

	err := tx.Raw(`
		SELECT request_hash, status, response_code, response_body, reservation_id::text AS reservation_id
		FROM idempotency_keys
		WHERE scope = ?
		  AND key = ?
		FOR UPDATE
	`, scope, key).Scan(&row).Error
	if err != nil {
		return CreateResult{}, err
	}

	if row.RequestHash != requestHash {
		return CreateResult{}, ErrIdempotencyKeyConflict
	}

	switch row.Status {
	case "in_progress":
		return CreateResult{}, ErrIdempotencyInProgress
	case "completed":
		result, err := createResultFromStoredBody(row.ResponseBody)
		if err != nil {
			return CreateResult{}, err
		}
		result.StatusCode = http.StatusOK
		result.IdempotentReplay = true
		return result, nil
	case "failed":
		return CreateResult{}, errorFromStoredBody(row.ResponseBody)
	default:
		return CreateResult{}, ErrIdempotencyInProgress
	}
}

func (s *PostgresStore) createFreshReservation(tx *gorm.DB, input CreateInput) (Reservation, error) {
	var updated struct {
		ID string
	}

	update := tx.Raw(`
		UPDATE items
		SET reserved_stock = reserved_stock + ?,
		    updated_at = NOW()
		WHERE id = ?::uuid
		  AND total_stock - reserved_stock >= ?
		RETURNING id::text AS id
	`, input.Quantity, input.ItemID, input.Quantity).Scan(&updated)
	if update.Error != nil {
		return Reservation{}, update.Error
	}
	if updated.ID == "" {
		exists, err := itemExists(tx, input.ItemID)
		if err != nil {
			return Reservation{}, err
		}
		if !exists {
			return Reservation{}, ErrItemNotFound
		}
		return Reservation{}, ErrInsufficientStock
	}

	reservation := Reservation{
		ID: uuid.NewString(),
	}

	if err := tx.Raw(`
		INSERT INTO reservations (id, item_id, session_id, quantity, status, expires_at, created_at, updated_at)
		VALUES (?::uuid, ?::uuid, ?, ?, 'active', NOW() + INTERVAL '60 seconds', NOW(), NOW())
		RETURNING
			id::text AS id,
			item_id::text AS item_id,
			session_id,
			quantity,
			status,
			expires_at,
			created_at
	`, reservation.ID, input.ItemID, input.SessionID, input.Quantity).Scan(&reservation).Error; err != nil {
		return Reservation{}, err
	}

	return reservation, nil
}

func expireActiveReservations(tx *gorm.DB) error {
	return tx.Exec(`
		WITH expired AS (
			UPDATE reservations
			SET status = 'expired',
			    expired_at = NOW(),
			    updated_at = NOW()
			WHERE status = 'active'
			  AND expires_at <= NOW()
			RETURNING item_id, quantity
		),
		totals AS (
			SELECT item_id, SUM(quantity) AS quantity
			FROM expired
			GROUP BY item_id
		)
		UPDATE items i
		SET reserved_stock = i.reserved_stock - totals.quantity,
		    updated_at = NOW()
		FROM totals
		WHERE i.id = totals.item_id
	`).Error
}

func lockReservation(tx *gorm.DB, sessionID string, reservationID string) (Reservation, error) {
	var reservation Reservation
	result := tx.Raw(`
		SELECT
			id::text AS id,
			item_id::text AS item_id,
			session_id,
			quantity,
			status,
			expires_at,
			created_at,
			released_at,
			expired_at
		FROM reservations
		WHERE id = ?::uuid
		  AND session_id = ?
		FOR UPDATE
	`, reservationID, sessionID).Scan(&reservation)
	if result.Error != nil {
		return Reservation{}, result.Error
	}
	if reservation.ID == "" {
		return Reservation{}, ErrReservationNotFound
	}
	return reservation, nil
}

func itemExists(tx *gorm.DB, itemID string) (bool, error) {
	var count int64
	err := tx.Raw(`SELECT COUNT(*) FROM items WHERE id = ?::uuid`, itemID).Scan(&count).Error
	return count > 0, err
}

func markIdempotencyCompleted(tx *gorm.DB, scope string, key string, reservation Reservation) error {
	body, _ := json.Marshal(createResponseBody{
		Reservation: reservationBody{
			ID:        reservation.ID,
			ItemID:    reservation.ItemID,
			Quantity:  reservation.Quantity,
			Status:    reservation.Status,
			ExpiresAt: reservation.ExpiresAt.UTC().Format(time.RFC3339),
			CreatedAt: reservation.CreatedAt.UTC().Format(time.RFC3339),
		},
	})
	return tx.Exec(`
		UPDATE idempotency_keys
		SET status = 'completed',
		    response_code = ?,
		    response_body = ?::jsonb,
		    reservation_id = ?::uuid,
		    updated_at = NOW()
		WHERE scope = ?
		  AND key = ?
	`, http.StatusCreated, string(body), reservation.ID, scope, key).Error
}

func markIdempotencyFailed(tx *gorm.DB, scope string, key string, statusCode int, code string) error {
	body, _ := json.Marshal(map[string]any{
		"error": map[string]any{
			"code": code,
		},
	})
	return tx.Exec(`
		UPDATE idempotency_keys
		SET status = 'failed',
		    response_code = ?,
		    response_body = ?::jsonb,
		    updated_at = NOW()
		WHERE scope = ?
		  AND key = ?
	`, statusCode, string(body), scope, key).Error
}

type createResponseBody struct {
	Reservation reservationBody `json:"reservation"`
}

type reservationBody struct {
	ID        string `json:"id"`
	ItemID    string `json:"itemId"`
	Quantity  int    `json:"quantity"`
	Status    string `json:"status"`
	ExpiresAt string `json:"expiresAt"`
	CreatedAt string `json:"createdAt"`
}

func createResultFromStoredBody(raw []byte) (CreateResult, error) {
	var body createResponseBody
	if err := json.Unmarshal(raw, &body); err != nil {
		return CreateResult{}, err
	}

	expiresAt, err := time.Parse(time.RFC3339, body.Reservation.ExpiresAt)
	if err != nil {
		return CreateResult{}, err
	}
	createdAt, err := time.Parse(time.RFC3339, body.Reservation.CreatedAt)
	if err != nil {
		return CreateResult{}, err
	}

	return CreateResult{
		Reservation: Reservation{
			ID:        body.Reservation.ID,
			ItemID:    body.Reservation.ItemID,
			Quantity:  body.Reservation.Quantity,
			Status:    body.Reservation.Status,
			ExpiresAt: expiresAt,
			CreatedAt: createdAt,
		},
	}, nil
}

func idempotentErrorPayload(err error) (int, string, bool) {
	switch {
	case errors.Is(err, ErrItemNotFound):
		return http.StatusNotFound, "item_not_found", true
	case errors.Is(err, ErrInsufficientStock):
		return http.StatusConflict, "insufficient_stock", true
	case errors.Is(err, ErrInvalidQuantity):
		return http.StatusBadRequest, "invalid_quantity", true
	default:
		return http.StatusInternalServerError, "internal_error", false
	}
}

func errorFromStoredBody(raw []byte) error {
	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return err
	}

	switch body.Error.Code {
	case "item_not_found":
		return ErrItemNotFound
	case "insufficient_stock":
		return ErrInsufficientStock
	case "invalid_quantity":
		return ErrInvalidQuantity
	default:
		return errors.New(body.Error.Code)
	}
}
