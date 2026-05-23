package reservations

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const testSessionID = "test-session"

func TestPostgresCreateReservationLastUnitConcurrency(t *testing.T) {
	db := openIntegrationDB(t)
	store := NewPostgresStore(db)
	service := NewService(store)
	itemID := seedItem(t, db, 1)

	successes, failures := runConcurrentCreates(t, service, itemID, 50, func(i int) string {
		return fmt.Sprintf("last-unit-%d", i)
	})

	if successes != 1 {
		t.Fatalf("expected exactly 1 success, got %d failures=%d", successes, failures)
	}
	assertStock(t, db, itemID, 1, 1, 0)
	assertReservationCount(t, db, itemID, StatusActive, 1)
}

func TestPostgresCreateReservationTenUnitsConcurrency(t *testing.T) {
	db := openIntegrationDB(t)
	store := NewPostgresStore(db)
	service := NewService(store)
	itemID := seedItem(t, db, 10)

	successes, failures := runConcurrentCreates(t, service, itemID, 100, func(i int) string {
		return fmt.Sprintf("ten-units-%d", i)
	})

	if successes != 10 || failures != 90 {
		t.Fatalf("expected 10 successes and 90 failures, got successes=%d failures=%d", successes, failures)
	}
	assertStock(t, db, itemID, 10, 10, 0)
	assertReservationCount(t, db, itemID, StatusActive, 10)
}

func TestPostgresCreateReservationParallelIdempotency(t *testing.T) {
	db := openIntegrationDB(t)
	store := NewPostgresStore(db)
	service := NewService(store)
	itemID := seedItem(t, db, 10)

	const workers = 20
	start := make(chan struct{})
	results := make(chan CreateResult, workers)
	errs := make(chan error, workers)

	for i := 0; i < workers; i++ {
		go func() {
			<-start
			result, err := service.CreateReservation(context.Background(), CreateInput{
				SessionID:      testSessionID,
				IdempotencyKey: "same-key",
				ItemID:         itemID,
				Quantity:       1,
			})
			if err != nil {
				errs <- err
				return
			}
			results <- result
		}()
	}
	close(start)

	var firstID string
	successes := 0
	for i := 0; i < workers; i++ {
		select {
		case err := <-errs:
			t.Fatalf("unexpected error: %v", err)
		case result := <-results:
			successes++
			if firstID == "" {
				firstID = result.Reservation.ID
			}
			if result.Reservation.ID != firstID {
				t.Fatalf("expected same reservation id %q, got %q", firstID, result.Reservation.ID)
			}
		}
	}

	if successes == 0 {
		t.Fatal("expected at least one successful idempotent result")
	}
	if successes != workers {
		t.Fatalf("expected %d successful idempotent outcomes, got %d", workers, successes)
	}
	assertStock(t, db, itemID, 10, 1, 9)
	assertReservationCount(t, db, itemID, StatusActive, 1)
}

func TestPostgresReleaseReservationDoubleCallReturnsStockOnce(t *testing.T) {
	db := openIntegrationDB(t)
	store := NewPostgresStore(db)
	service := NewService(store)
	itemID := seedItem(t, db, 1)

	created, err := service.CreateReservation(context.Background(), CreateInput{
		SessionID:      testSessionID,
		IdempotencyKey: "release-once",
		ItemID:         itemID,
		Quantity:       1,
	})
	if err != nil {
		t.Fatalf("create reservation: %v", err)
	}

	const workers = 2
	start := make(chan struct{})
	results := make(chan ReleaseResult, workers)
	errs := make(chan error, workers)

	for i := 0; i < workers; i++ {
		go func() {
			<-start
			result, err := service.ReleaseReservation(context.Background(), testSessionID, created.Reservation.ID)
			if err != nil {
				errs <- err
				return
			}
			results <- result
		}()
	}
	close(start)

	stockReturned := 0
	for i := 0; i < workers; i++ {
		select {
		case err := <-errs:
			t.Fatalf("unexpected release error: %v", err)
		case result := <-results:
			if result.StockReturned {
				stockReturned++
			}
		}
	}

	if stockReturned != 1 {
		t.Fatalf("expected stock returned once, got %d", stockReturned)
	}
	assertStock(t, db, itemID, 1, 0, 1)
	assertReservationCount(t, db, itemID, StatusReleased, 1)
}

func TestPostgresIdempotencyConflictDoesNotChangeStock(t *testing.T) {
	db := openIntegrationDB(t)
	store := NewPostgresStore(db)
	service := NewService(store)
	itemID := seedItem(t, db, 10)

	_, err := service.CreateReservation(context.Background(), CreateInput{
		SessionID:      testSessionID,
		IdempotencyKey: "conflict-key",
		ItemID:         itemID,
		Quantity:       1,
	})
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	_, err = service.CreateReservation(context.Background(), CreateInput{
		SessionID:      testSessionID,
		IdempotencyKey: "conflict-key",
		ItemID:         itemID,
		Quantity:       2,
	})
	if !errors.Is(err, ErrIdempotencyKeyConflict) {
		t.Fatalf("expected idempotency conflict, got %v", err)
	}

	assertStock(t, db, itemID, 10, 1, 9)
	assertReservationCount(t, db, itemID, StatusActive, 1)
}

func TestPostgresExpirationTwiceReturnsStockOnce(t *testing.T) {
	db := openIntegrationDB(t)
	store := NewPostgresStore(db)
	service := NewService(store)
	itemID := seedItem(t, db, 1)

	created, err := service.CreateReservation(context.Background(), CreateInput{
		SessionID:      testSessionID,
		IdempotencyKey: "expire-once",
		ItemID:         itemID,
		Quantity:       1,
	})
	if err != nil {
		t.Fatalf("create reservation: %v", err)
	}
	if err := db.Exec(`
		UPDATE reservations
		SET created_at = NOW() - INTERVAL '2 minutes',
		    expires_at = NOW() - INTERVAL '1 minute'
		WHERE id = ?::uuid
	`, created.Reservation.ID).Error; err != nil {
		t.Fatalf("force expiration: %v", err)
	}

	if err := store.ExpireActiveReservations(context.Background()); err != nil {
		t.Fatalf("first expiration: %v", err)
	}
	if err := store.ExpireActiveReservations(context.Background()); err != nil {
		t.Fatalf("second expiration: %v", err)
	}

	assertStock(t, db, itemID, 1, 0, 1)
	assertReservationCount(t, db, itemID, StatusExpired, 1)
}

func TestPostgresReleaseVsExpirationReturnsStockOnce(t *testing.T) {
	db := openIntegrationDB(t)
	store := NewPostgresStore(db)
	service := NewService(store)
	itemID := seedItem(t, db, 1)

	created, err := service.CreateReservation(context.Background(), CreateInput{
		SessionID:      testSessionID,
		IdempotencyKey: "release-vs-expire",
		ItemID:         itemID,
		Quantity:       1,
	})
	if err != nil {
		t.Fatalf("create reservation: %v", err)
	}
	if err := db.Exec(`
		UPDATE reservations
		SET created_at = NOW() - INTERVAL '2 minutes',
		    expires_at = NOW() - INTERVAL '1 minute'
		WHERE id = ?::uuid
	`, created.Reservation.ID).Error; err != nil {
		t.Fatalf("force expiration: %v", err)
	}

	start := make(chan struct{})
	errs := make(chan error, 2)
	go func() {
		<-start
		errs <- store.ExpireActiveReservations(context.Background())
	}()
	go func() {
		<-start
		_, err := service.ReleaseReservation(context.Background(), testSessionID, created.Reservation.ID)
		if errors.Is(err, ErrReservationNotFound) {
			err = nil
		}
		errs <- err
	}()
	close(start)

	for i := 0; i < 2; i++ {
		if err := <-errs; err != nil {
			t.Fatalf("unexpected race error: %v", err)
		}
	}

	assertStock(t, db, itemID, 1, 0, 1)
	assertTerminalReservation(t, db, itemID)
}

func openIntegrationDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open integration db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("unwrap integration db: %v", err)
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)

	resetDatabase(t, db)
	return db
}

func resetDatabase(t *testing.T, db *gorm.DB) {
	t.Helper()

	if err := db.Exec(`
		DROP TABLE IF EXISTS idempotency_keys;
		DROP TABLE IF EXISTS reservations;
		DROP TABLE IF EXISTS items;
	`).Error; err != nil {
		t.Fatalf("drop schema: %v", err)
	}

	migrationPath := filepath.Join("..", "..", "..", "db", "migrations", "001_initial_schema.sql")
	raw, err := os.ReadFile(migrationPath)
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	if err := db.Exec(string(raw)).Error; err != nil {
		t.Fatalf("apply migration: %v", err)
	}
}

func seedItem(t *testing.T, db *gorm.DB, stock int) string {
	t.Helper()

	itemID := fmt.Sprintf("10000000-0000-0000-0000-%012d", stock)
	if err := db.Exec(`
		INSERT INTO items (id, name, total_stock, reserved_stock)
		VALUES (?::uuid, ?, ?, 0)
	`, itemID, fmt.Sprintf("Test Item %d", stock), stock).Error; err != nil {
		t.Fatalf("seed item: %v", err)
	}
	return itemID
}

func runConcurrentCreates(t *testing.T, service *ReservationService, itemID string, workers int, key func(int) string) (int, int) {
	t.Helper()

	start := make(chan struct{})
	var wg sync.WaitGroup
	var mu sync.Mutex
	successes := 0
	failures := 0

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_, err := service.CreateReservation(context.Background(), CreateInput{
				SessionID:      testSessionID,
				IdempotencyKey: key(i),
				ItemID:         itemID,
				Quantity:       1,
			})

			mu.Lock()
			defer mu.Unlock()
			if err == nil {
				successes++
				return
			}
			if errors.Is(err, ErrInsufficientStock) || errors.Is(err, ErrIdempotencyInProgress) {
				failures++
				return
			}
			t.Errorf("unexpected create error: %v", err)
		}(i)
	}

	close(start)
	wg.Wait()

	return successes, failures
}

func assertStock(t *testing.T, db *gorm.DB, itemID string, total int, reserved int, available int) {
	t.Helper()

	var row struct {
		TotalStock     int
		ReservedStock  int
		AvailableStock int
	}
	if err := db.Raw(`
		SELECT total_stock, reserved_stock, total_stock - reserved_stock AS available_stock
		FROM items
		WHERE id = ?::uuid
	`, itemID).Scan(&row).Error; err != nil {
		t.Fatalf("read stock: %v", err)
	}
	if row.TotalStock != total || row.ReservedStock != reserved || row.AvailableStock != available {
		t.Fatalf("expected stock total=%d reserved=%d available=%d, got total=%d reserved=%d available=%d",
			total, reserved, available, row.TotalStock, row.ReservedStock, row.AvailableStock)
	}
}

func assertReservationCount(t *testing.T, db *gorm.DB, itemID string, status string, expected int64) {
	t.Helper()

	var count int64
	if err := db.Raw(`
		SELECT COUNT(*)
		FROM reservations
		WHERE item_id = ?::uuid
		  AND status = ?
	`, itemID, status).Scan(&count).Error; err != nil {
		t.Fatalf("count reservations: %v", err)
	}
	if count != expected {
		t.Fatalf("expected %d %s reservations, got %d", expected, status, count)
	}
}

func assertTerminalReservation(t *testing.T, db *gorm.DB, itemID string) {
	t.Helper()

	var row struct {
		Status string
		Count  int64
	}
	if err := db.Raw(`
		SELECT status, COUNT(*) AS count
		FROM reservations
		WHERE item_id = ?::uuid
		GROUP BY status
	`, itemID).Scan(&row).Error; err != nil {
		t.Fatalf("read terminal reservation: %v", err)
	}
	if row.Count != 1 {
		t.Fatalf("expected one reservation, got %d", row.Count)
	}
	if row.Status != StatusExpired && row.Status != StatusReleased {
		t.Fatalf("expected terminal reservation, got status %q", row.Status)
	}
}
