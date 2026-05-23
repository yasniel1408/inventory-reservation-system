package reservations

import (
	"context"
	"testing"
)

type fakeStore struct {
	expireCalls  int
	createCalls  int
	listCalls    int
	releaseCalls int
	lastHash     string
}

func (f *fakeStore) ExpireActiveReservations(ctx context.Context) error {
	f.expireCalls++
	return nil
}

func (f *fakeStore) CreateReservation(ctx context.Context, input CreateInput, requestHash string) (CreateResult, error) {
	f.createCalls++
	f.lastHash = requestHash
	return CreateResult{Reservation: Reservation{ID: "10000000-0000-0000-0000-000000000001"}}, nil
}

func (f *fakeStore) ListActiveReservations(ctx context.Context, sessionID string) ([]Reservation, error) {
	f.listCalls++
	return []Reservation{{ID: "10000000-0000-0000-0000-000000000001", SessionID: sessionID}}, nil
}

func (f *fakeStore) ReleaseReservation(ctx context.Context, sessionID string, reservationID string) (ReleaseResult, error) {
	f.releaseCalls++
	return ReleaseResult{Reservation: Reservation{ID: reservationID, SessionID: sessionID}}, nil
}

func TestCreateReservationRejectsInvalidQuantityBeforeStore(t *testing.T) {
	store := &fakeStore{}
	service := NewService(store)

	_, err := service.CreateReservation(context.Background(), CreateInput{
		SessionID:      "session-1",
		IdempotencyKey: "key-1",
		ItemID:         "00000000-0000-0000-0000-000000000001",
		Quantity:       0,
	})

	if err != ErrInvalidQuantity {
		t.Fatalf("expected invalid quantity, got %v", err)
	}
	if store.expireCalls != 0 || store.createCalls != 0 {
		t.Fatalf("expected no store calls, got expire=%d create=%d", store.expireCalls, store.createCalls)
	}
}

func TestCreateReservationExpiresBeforeStoreCreateAndHashesRequest(t *testing.T) {
	store := &fakeStore{}
	service := NewService(store)

	_, err := service.CreateReservation(context.Background(), CreateInput{
		SessionID:      "session-1",
		IdempotencyKey: "key-1",
		ItemID:         "00000000-0000-0000-0000-000000000001",
		Quantity:       2,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.expireCalls != 1 || store.createCalls != 1 {
		t.Fatalf("expected expire and create once, got expire=%d create=%d", store.expireCalls, store.createCalls)
	}
	if store.lastHash == "" {
		t.Fatal("expected canonical request hash")
	}
}

func TestCreateReservationRejectsMalformedItemIDBeforeStore(t *testing.T) {
	store := &fakeStore{}
	service := NewService(store)

	_, err := service.CreateReservation(context.Background(), CreateInput{
		SessionID:      "session-1",
		IdempotencyKey: "key-1",
		ItemID:         "not-a-uuid",
		Quantity:       1,
	})

	if err != ErrInvalidItemID {
		t.Fatalf("expected invalid item id, got %v", err)
	}
	if store.expireCalls != 0 || store.createCalls != 0 {
		t.Fatalf("expected no store calls, got expire=%d create=%d", store.expireCalls, store.createCalls)
	}
}

func TestReleaseReservationRejectsMalformedReservationIDBeforeStore(t *testing.T) {
	store := &fakeStore{}
	service := NewService(store)

	_, err := service.ReleaseReservation(context.Background(), "session-1", "not-a-uuid")

	if err != ErrInvalidReservationID {
		t.Fatalf("expected invalid reservation id, got %v", err)
	}
	if store.expireCalls != 0 || store.releaseCalls != 0 {
		t.Fatalf("expected no store calls, got expire=%d release=%d", store.expireCalls, store.releaseCalls)
	}
}

func TestListActiveReservationsExpiresBeforeListing(t *testing.T) {
	store := &fakeStore{}
	service := NewService(store)

	_, err := service.ListActiveReservations(context.Background(), "session-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.expireCalls != 1 || store.listCalls != 1 {
		t.Fatalf("expected expire and list once, got expire=%d list=%d", store.expireCalls, store.listCalls)
	}
}

func TestReleaseReservationExpiresBeforeRelease(t *testing.T) {
	store := &fakeStore{}
	service := NewService(store)

	_, err := service.ReleaseReservation(context.Background(), "session-1", "10000000-0000-0000-0000-000000000001")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.expireCalls != 1 || store.releaseCalls != 1 {
		t.Fatalf("expected expire and release once, got expire=%d release=%d", store.expireCalls, store.releaseCalls)
	}
}
