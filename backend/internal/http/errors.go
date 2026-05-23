package httpapi

import (
	"errors"
	"net/http"

	"inventory-reservation-system/backend/internal/reservations"
)

type apiError struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}

func writeError(c responseWriter, status int, code string, message string) {
	c.JSON(status, apiError{
		Error: errorBody{
			Code:    code,
			Message: message,
			Details: map[string]any{},
		},
	})
}

func writeDomainError(c responseWriter, err error) {
	switch {
	case errors.Is(err, reservations.ErrInvalidQuantity):
		writeError(c, http.StatusBadRequest, "invalid_quantity", "Quantity must be a positive integer.")
	case errors.Is(err, reservations.ErrInvalidIdempotencyKey):
		writeError(c, http.StatusBadRequest, "missing_idempotency_key", "Idempotency-Key header is required.")
	case errors.Is(err, reservations.ErrInvalidItemID):
		writeError(c, http.StatusBadRequest, "invalid_request", "itemId must be a valid UUID.")
	case errors.Is(err, reservations.ErrInvalidReservationID):
		writeError(c, http.StatusBadRequest, "invalid_request", "reservation id must be a valid UUID.")
	case errors.Is(err, reservations.ErrItemNotFound):
		writeError(c, http.StatusNotFound, "item_not_found", "Item was not found.")
	case errors.Is(err, reservations.ErrInsufficientStock):
		writeError(c, http.StatusConflict, "insufficient_stock", "Not enough stock available.")
	case errors.Is(err, reservations.ErrIdempotencyKeyConflict):
		writeError(c, http.StatusConflict, "idempotency_key_conflict", "Idempotency key was already used with a different payload.")
	case errors.Is(err, reservations.ErrIdempotencyInProgress):
		writeError(c, http.StatusTooEarly, "idempotency_in_progress", "A request with this idempotency key is still in progress.")
	case errors.Is(err, reservations.ErrReservationNotFound):
		writeError(c, http.StatusNotFound, "reservation_not_found", "Reservation was not found.")
	default:
		writeError(c, http.StatusInternalServerError, "internal_error", "Unexpected internal error.")
	}
}
