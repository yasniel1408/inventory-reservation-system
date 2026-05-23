package httpapi

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"inventory-reservation-system/backend/internal/items"
	"inventory-reservation-system/backend/internal/reservations"
)

const defaultSessionID = "anonymous-session"

type Dependencies struct {
	Items        items.Service
	Reservations reservations.Service
}

type responseWriter interface {
	JSON(code int, obj any)
}

type createReservationRequest struct {
	ItemID   string `json:"itemId"`
	Quantity int    `json:"quantity"`
}

func NewRouter(deps Dependencies) http.Handler {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/items", func(c *gin.Context) {
		result, err := deps.Items.ListItems(c.Request.Context())
		if err != nil {
			writeDomainError(c, err)
			return
		}

		itemsResponse := make([]gin.H, 0, len(result))
		for _, item := range result {
			itemsResponse = append(itemsResponse, gin.H{
				"id":             item.ID,
				"name":           item.Name,
				"totalStock":     item.TotalStock,
				"reservedStock":  item.ReservedStock,
				"availableStock": item.AvailableStock(),
			})
		}

		c.JSON(http.StatusOK, gin.H{"items": itemsResponse})
	})

	router.POST("/reservations", func(c *gin.Context) {
		idempotencyKey := strings.TrimSpace(c.GetHeader("Idempotency-Key"))
		if idempotencyKey == "" {
			writeError(c, http.StatusBadRequest, "missing_idempotency_key", "Idempotency-Key header is required.")
			return
		}

		var req createReservationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			writeError(c, http.StatusBadRequest, "invalid_request", "Request body is invalid.")
			return
		}

		result, err := deps.Reservations.CreateReservation(c.Request.Context(), reservations.CreateInput{
			SessionID:      sessionID(c),
			IdempotencyKey: idempotencyKey,
			ItemID:         strings.TrimSpace(req.ItemID),
			Quantity:       req.Quantity,
		})
		if err != nil {
			writeDomainError(c, err)
			return
		}

		status := result.StatusCode
		if status == 0 {
			status = http.StatusCreated
		}

		response := gin.H{"reservation": reservationResponse(result.Reservation)}
		if result.IdempotentReplay {
			response["idempotentReplay"] = true
		}

		c.JSON(status, response)
	})

	router.GET("/reservations", func(c *gin.Context) {
		result, err := deps.Reservations.ListActiveReservations(c.Request.Context(), sessionID(c))
		if err != nil {
			writeDomainError(c, err)
			return
		}

		response := make([]gin.H, 0, len(result))
		for _, reservation := range result {
			response = append(response, reservationResponse(reservation))
		}

		c.JSON(http.StatusOK, gin.H{"reservations": response})
	})

	router.DELETE("/reservations/:id", func(c *gin.Context) {
		result, err := deps.Reservations.ReleaseReservation(c.Request.Context(), sessionID(c), c.Param("id"))
		if err != nil {
			writeDomainError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"reservation":   reservationResponse(result.Reservation),
			"stockReturned": result.StockReturned,
			"noop":          result.Noop,
		})
	})

	return router
}

func sessionID(c *gin.Context) string {
	value := strings.TrimSpace(c.GetHeader("X-Session-ID"))
	if value == "" {
		return defaultSessionID
	}
	return value
}

func reservationResponse(reservation reservations.Reservation) gin.H {
	body := gin.H{
		"id":        reservation.ID,
		"itemId":    reservation.ItemID,
		"quantity":  reservation.Quantity,
		"status":    reservation.Status,
		"expiresAt": reservation.ExpiresAt.UTC().Format(timeFormat),
		"createdAt": reservation.CreatedAt.UTC().Format(timeFormat),
	}

	if reservation.ItemName != "" {
		body["itemName"] = reservation.ItemName
	}
	if reservation.ReleasedAt != nil {
		body["releasedAt"] = reservation.ReleasedAt.UTC().Format(timeFormat)
	}
	if reservation.ExpiredAt != nil {
		body["expiredAt"] = reservation.ExpiredAt.UTC().Format(timeFormat)
	}

	return body
}

const timeFormat = "2006-01-02T15:04:05Z07:00"
