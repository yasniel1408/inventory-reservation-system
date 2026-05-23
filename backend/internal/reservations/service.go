package reservations

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

type Store interface {
	ExpireActiveReservations(ctx context.Context) error
	CreateReservation(ctx context.Context, input CreateInput, requestHash string) (CreateResult, error)
	ListActiveReservations(ctx context.Context, sessionID string) ([]Reservation, error)
	ReleaseReservation(ctx context.Context, sessionID string, reservationID string) (ReleaseResult, error)
}

type ReservationService struct {
	store Store
}

func NewService(store Store) *ReservationService {
	return &ReservationService{store: store}
}

func (s *ReservationService) CreateReservation(ctx context.Context, input CreateInput) (CreateResult, error) {
	input.SessionID = normalizeSession(input.SessionID)
	input.ItemID = strings.TrimSpace(input.ItemID)
	input.IdempotencyKey = strings.TrimSpace(input.IdempotencyKey)

	if input.IdempotencyKey == "" {
		return CreateResult{}, ErrInvalidIdempotencyKey
	}
	if input.ItemID == "" || input.Quantity <= 0 {
		return CreateResult{}, ErrInvalidQuantity
	}
	if _, err := uuid.Parse(input.ItemID); err != nil {
		return CreateResult{}, ErrInvalidItemID
	}

	if err := s.store.ExpireActiveReservations(ctx); err != nil {
		return CreateResult{}, err
	}

	return s.store.CreateReservation(ctx, input, canonicalRequestHash(input))
}

func (s *ReservationService) ListActiveReservations(ctx context.Context, sessionID string) ([]Reservation, error) {
	sessionID = normalizeSession(sessionID)
	if err := s.store.ExpireActiveReservations(ctx); err != nil {
		return nil, err
	}

	return s.store.ListActiveReservations(ctx, sessionID)
}

func (s *ReservationService) ReleaseReservation(ctx context.Context, sessionID string, reservationID string) (ReleaseResult, error) {
	sessionID = normalizeSession(sessionID)
	reservationID = strings.TrimSpace(reservationID)
	if reservationID == "" {
		return ReleaseResult{}, ErrReservationNotFound
	}
	if _, err := uuid.Parse(reservationID); err != nil {
		return ReleaseResult{}, ErrInvalidReservationID
	}
	if err := s.store.ExpireActiveReservations(ctx); err != nil {
		return ReleaseResult{}, err
	}

	return s.store.ReleaseReservation(ctx, sessionID, reservationID)
}

func canonicalRequestHash(input CreateInput) string {
	payload := struct {
		Method    string `json:"method"`
		Path      string `json:"path"`
		SessionID string `json:"sessionId"`
		ItemID    string `json:"itemId"`
		Quantity  int    `json:"quantity"`
	}{
		Method:    "POST",
		Path:      "/reservations",
		SessionID: normalizeSession(input.SessionID),
		ItemID:    strings.TrimSpace(input.ItemID),
		Quantity:  input.Quantity,
	}

	raw, _ := json.Marshal(payload)
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func normalizeSession(sessionID string) string {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return "anonymous-session"
	}
	return sessionID
}
