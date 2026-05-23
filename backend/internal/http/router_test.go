package httpapi_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"inventory-reservation-system/backend/internal/http"
	"inventory-reservation-system/backend/internal/items"
	"inventory-reservation-system/backend/internal/reservations"
)

type fakeItemsService struct {
	items []items.Item
	err   error
}

func (f fakeItemsService) ListItems(ctx context.Context) ([]items.Item, error) {
	return f.items, f.err
}

type fakeReservationsService struct {
	createResult reservations.CreateResult
	createErr    error
	listResult   []reservations.Reservation
	listErr      error
	releaseRes   reservations.ReleaseResult
	releaseErr   error
}

func (f fakeReservationsService) CreateReservation(ctx context.Context, input reservations.CreateInput) (reservations.CreateResult, error) {
	return f.createResult, f.createErr
}

func (f fakeReservationsService) ListActiveReservations(ctx context.Context, sessionID string) ([]reservations.Reservation, error) {
	return f.listResult, f.listErr
}

func (f fakeReservationsService) ReleaseReservation(ctx context.Context, sessionID string, reservationID string) (reservations.ReleaseResult, error) {
	return f.releaseRes, f.releaseErr
}

func TestGetItemsReturnsInventoryShape(t *testing.T) {
	router := httpapi.NewRouter(httpapi.Dependencies{
		Items: fakeItemsService{items: []items.Item{{
			ID:            "00000000-0000-0000-0000-000000000001",
			Name:          "Standard Widget",
			TotalStock:    100,
			ReservedStock: 3,
		}}},
		Reservations: fakeReservationsService{},
	})

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", res.Code, res.Body.String())
	}

	var body struct {
		Items []struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			TotalStock     int    `json:"totalStock"`
			ReservedStock  int    `json:"reservedStock"`
			AvailableStock int    `json:"availableStock"`
		} `json:"items"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(body.Items) != 1 {
		t.Fatalf("expected one item, got %d", len(body.Items))
	}
	if body.Items[0].AvailableStock != 97 {
		t.Fatalf("expected available stock 97, got %d", body.Items[0].AvailableStock)
	}
}

func TestPostReservationsRequiresIdempotencyKey(t *testing.T) {
	router := httpapi.NewRouter(httpapi.Dependencies{
		Items:        fakeItemsService{},
		Reservations: fakeReservationsService{},
	})

	req := httptest.NewRequest(http.MethodPost, "/reservations", strings.NewReader(`{"itemId":"item-1","quantity":1}`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", res.Code, res.Body.String())
	}
	assertErrorCode(t, res.Body.Bytes(), "missing_idempotency_key")
}

func TestPostReservationsMapsDomainConflict(t *testing.T) {
	router := httpapi.NewRouter(httpapi.Dependencies{
		Items: fakeItemsService{},
		Reservations: fakeReservationsService{
			createErr: reservations.ErrInsufficientStock,
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/reservations", strings.NewReader(`{"itemId":"item-1","quantity":2}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "reserve-1")
	req.Header.Set("X-Session-ID", "session-1")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", res.Code, res.Body.String())
	}
	assertErrorCode(t, res.Body.Bytes(), "insufficient_stock")
}

func TestPostReservationsReturnsCreatedReservation(t *testing.T) {
	expiresAt := time.Date(2026, 5, 22, 12, 1, 0, 0, time.UTC)
	createdAt := expiresAt.Add(-time.Minute)
	router := httpapi.NewRouter(httpapi.Dependencies{
		Items: fakeItemsService{},
		Reservations: fakeReservationsService{
			createResult: reservations.CreateResult{
				StatusCode: http.StatusCreated,
				Reservation: reservations.Reservation{
					ID:        "res-1",
					ItemID:    "item-1",
					Quantity:  2,
					Status:    reservations.StatusActive,
					ExpiresAt: expiresAt,
					CreatedAt: createdAt,
				},
			},
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/reservations", strings.NewReader(`{"itemId":"item-1","quantity":2}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "reserve-1")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", res.Code, res.Body.String())
	}

	var body struct {
		Reservation struct {
			ID        string `json:"id"`
			ItemID    string `json:"itemId"`
			Quantity  int    `json:"quantity"`
			Status    string `json:"status"`
			ExpiresAt string `json:"expiresAt"`
			CreatedAt string `json:"createdAt"`
		} `json:"reservation"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body.Reservation.ID != "res-1" || body.Reservation.Status != "active" {
		t.Fatalf("unexpected reservation body: %+v", body.Reservation)
	}
}

func TestDeleteReservationReturnsNoopRelease(t *testing.T) {
	router := httpapi.NewRouter(httpapi.Dependencies{
		Items: fakeItemsService{},
		Reservations: fakeReservationsService{
			releaseRes: reservations.ReleaseResult{
				Reservation: reservations.Reservation{
					ID:     "res-1",
					Status: reservations.StatusReleased,
				},
				StockReturned: false,
				Noop:          true,
			},
		},
	})

	req := httptest.NewRequest(http.MethodDelete, "/reservations/res-1", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", res.Code, res.Body.String())
	}

	var body struct {
		StockReturned bool `json:"stockReturned"`
		Noop          bool `json:"noop"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if !body.Noop || body.StockReturned {
		t.Fatalf("expected noop release without stock return, got %+v", body)
	}
}

func TestInternalErrorsUseStableShape(t *testing.T) {
	router := httpapi.NewRouter(httpapi.Dependencies{
		Items:        fakeItemsService{err: errors.New("db down")},
		Reservations: fakeReservationsService{},
	})

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", res.Code, res.Body.String())
	}
	assertErrorCode(t, res.Body.Bytes(), "internal_error")
}

func TestCorsPreflightAllowsFrontendReservationHeaders(t *testing.T) {
	router := httpapi.NewRouter(httpapi.Dependencies{
		Items:        fakeItemsService{},
		Reservations: fakeReservationsService{},
	})

	req := httptest.NewRequest(http.MethodOptions, "/reservations", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	req.Header.Set("Access-Control-Request-Headers", "content-type,idempotency-key")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", res.Code, res.Body.String())
	}
	if got := res.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("expected wildcard allow origin, got %q", got)
	}
	if got := res.Header().Get("Access-Control-Allow-Headers"); !strings.Contains(got, "Idempotency-Key") {
		t.Fatalf("expected idempotency header to be allowed, got %q", got)
	}
}

func assertErrorCode(t *testing.T, raw []byte, expected string) {
	t.Helper()

	var body struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("invalid error json: %v", err)
	}
	if body.Error.Code != expected {
		t.Fatalf("expected error code %q, got %q", expected, body.Error.Code)
	}
}
