package repo

import (
	"context"

	"github.com/kimvlry/simple-order-service/internal/domain"
)

type OrderRepo interface {
	GetById(ctx context.Context, uid string) (*domain.Order, error)
	GetAll(ctx context.Context) ([]domain.Order, error)
	Save(ctx context.Context, order *domain.Order) error
}
