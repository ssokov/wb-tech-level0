package service

import (
	"context"
	"github.com/kimvlry/simple-order-service/internal/domain"
	"github.com/kimvlry/simple-order-service/internal/repo"
	"github.com/kimvlry/simple-order-service/internal/service/redis"
)

type OrderService struct {
	repo        repo.OrderRepo
	redisClient *redis.Client
}

func NewOrderService(repo repo.OrderRepo, client *redis.Client) *OrderService {
	return &OrderService{
		repo:        repo,
		redisClient: client,
	}
}

func (s *OrderService) GetOrderByID(id string, ctx context.Context) (*domain.Order, error) {
	order, err := s.redisClient.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return order, nil
	}

	if order, err = s.repo.GetById(ctx, id); err != nil {
		return nil, err
	}
	if order == nil {
		_ = s.redisClient.SaveOrder(ctx, order)
	}
	return order, nil
}
