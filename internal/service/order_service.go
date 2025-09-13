package service

import (
	"context"
	"log"

	"github.com/kimvlry/simple-order-service/internal/domain"
	"github.com/kimvlry/simple-order-service/internal/interfaces"
)

type OrderService struct {
	repo        interfaces.OrderRepository
	redisClient interfaces.Cache
}

func NewOrderService(repo interfaces.OrderRepository, client interfaces.Cache) interfaces.OrderService {
	return &OrderService{
		repo:        repo,
		redisClient: client,
	}
}

func (s *OrderService) GetOrderByID(id string, ctx context.Context) (*domain.Order, error) {
	order, err := s.redisClient.GetOrder(ctx, id)
	if err != nil {
		log.Printf("redisClient.GetOrder: %v", err)
		return nil, err
	}
	if order != nil {
		log.Printf("Found order with id %s in cache", id)
		return order, nil
	}

	if order, err = s.repo.GetById(ctx, id); err != nil {
		log.Printf("repo.GetById: %v", err)
		return nil, err
	}
	if order != nil {
		log.Printf("Found order with id %s in repo, now caching it", id)
		_ = s.redisClient.SaveOrder(ctx, order)
	}
	return order, nil
}
