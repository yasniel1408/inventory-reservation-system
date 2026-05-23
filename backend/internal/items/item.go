package items

import (
	"context"

	"gorm.io/gorm"
)

type Item struct {
	ID            string
	Name          string
	TotalStock    int
	ReservedStock int
}

func (i Item) AvailableStock() int {
	return i.TotalStock - i.ReservedStock
}

type Service interface {
	ListItems(ctx context.Context) ([]Item, error)
}

type Repository interface {
	ListItems(ctx context.Context) ([]Item, error)
}

type Expirer interface {
	ExpireActiveReservations(ctx context.Context) error
}

type ItemService struct {
	repo    Repository
	expirer Expirer
}

func NewService(repo Repository, expirer Expirer) *ItemService {
	return &ItemService{repo: repo, expirer: expirer}
}

func (s *ItemService) ListItems(ctx context.Context) ([]Item, error) {
	if s.expirer != nil {
		if err := s.expirer.ExpireActiveReservations(ctx); err != nil {
			return nil, err
		}
	}

	return s.repo.ListItems(ctx)
}

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ListItems(ctx context.Context) ([]Item, error) {
	var result []Item
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			id::text AS id,
			name,
			total_stock,
			reserved_stock
		FROM items
		ORDER BY name ASC
	`).Scan(&result).Error
	return result, err
}
