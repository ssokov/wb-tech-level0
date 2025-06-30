package service

import (
	"context"
	"github.com/kimvlry/simple-order-service/internal/domain"
	"github.com/kimvlry/simple-order-service/internal/repo"
	"github.com/kimvlry/simple-order-service/internal/service/redis"
)

type OrderService struct {
	repo        repo.OrderRepo
	redisClient *redis.RedisClient
	ctx         context.Context
}

func NewOrderService(repo repo.OrderRepo, client *redis.RedisClient) *OrderService {
	return &OrderService{
		repo:        repo,
		ctx:         context.Background(),
		redisClient: client,
	}
}

func (s *OrderService) GetOrderByID(id string) (*domain.Order, error) {
	order, err := s.redisClient.GetOrder(s.ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return order, nil
	}

	if order, err = s.repo.GetById(s.ctx, id); err != nil {
		return nil, err
	}
	if order == nil {
		_ = s.redisClient.SaveOrder(s.ctx, order)
	}
	return order, nil
}
