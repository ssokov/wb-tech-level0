package interfaces

import (
	"context"

	"github.com/kimvlry/simple-order-service/internal/domain"
)

// OrderRepository интерфейс для работы с заказами в БД
type OrderRepository interface {
	GetById(ctx context.Context, id string) (*domain.Order, error)
	GetAll(ctx context.Context) ([]domain.Order, error)
	Save(ctx context.Context, order *domain.Order) error
	Close() error
}

// Cache интерфейс для работы с кэшем
type Cache interface {
	GetOrder(ctx context.Context, orderUid string) (*domain.Order, error)
	SaveOrder(ctx context.Context, order *domain.Order) error
	RestoreCache(ctx context.Context) error
	Close() error
}

// MessageConsumer интерфейс для работы с сообщениями
type MessageConsumer interface {
	Consume(ctx context.Context) error
	Close() error
}

// OrderService интерфейс для бизнес-логики заказов
type OrderService interface {
	GetOrderByID(id string, ctx context.Context) (*domain.Order, error)
}
