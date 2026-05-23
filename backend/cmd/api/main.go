package main

import (
	"log"
	"net/http"

	"inventory-reservation-system/backend/internal/config"
	httpapi "inventory-reservation-system/backend/internal/http"
	"inventory-reservation-system/backend/internal/items"
	"inventory-reservation-system/backend/internal/reservations"
	"inventory-reservation-system/backend/internal/store"
)

func main() {
	cfg := config.Load()

	db, err := store.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}

	reservationStore := reservations.NewPostgresStore(db)
	itemService := items.NewService(items.NewPostgresRepository(db), reservationStore)
	reservationService := reservations.NewService(reservationStore)

	router := httpapi.NewRouter(httpapi.Dependencies{
		Items:        itemService,
		Reservations: reservationService,
	})

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	log.Printf("backend listening on :%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %v", err)
	}
}
