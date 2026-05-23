package reservations

import (
	"context"
	"errors"
	"time"
)

const (
	StatusActive    = "active"
	StatusReleased  = "released"
	StatusExpired   = "expired"
	StatusConfirmed = "confirmed"
)

var (
	ErrInvalidQuantity         = errors.New("invalid quantity")
	ErrItemNotFound            = errors.New("item not found")
	ErrInsufficientStock       = errors.New("insufficient stock")
	ErrIdempotencyKeyConflict  = errors.New("idempotency key conflict")
	ErrIdempotencyInProgress   = errors.New("idempotency in progress")
	ErrReservationNotFound     = errors.New("reservation not found")
	ErrInvalidIdempotencyKey   = errors.New("invalid idempotency key")
	ErrInvalidReservationOwner = errors.New("invalid reservation owner")
	ErrInvalidItemID           = errors.New("invalid item id")
	ErrInvalidReservationID    = errors.New("invalid reservation id")
)

type Reservation struct {
	ID         string
	ItemID     string
	ItemName   string
	SessionID  string
	Quantity   int
	Status     string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	ReleasedAt *time.Time
	ExpiredAt  *time.Time
}

type CreateInput struct {
	SessionID      string
	IdempotencyKey string
	ItemID         string
	Quantity       int
}

type CreateResult struct {
	Reservation      Reservation
	StatusCode       int
	IdempotentReplay bool
}

type ReleaseResult struct {
	Reservation   Reservation
	StockReturned bool
	Noop          bool
}

type Service interface {
	CreateReservation(ctx context.Context, input CreateInput) (CreateResult, error)
	ListActiveReservations(ctx context.Context, sessionID string) ([]Reservation, error)
	ReleaseReservation(ctx context.Context, sessionID string, reservationID string) (ReleaseResult, error)
}
