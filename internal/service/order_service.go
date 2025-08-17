package service

import (
	"context"
	"simple-order-service/internal/domain"
	"simple-order-service/internal/repo"
	"simple-order-service/internal/service/redis"

	"log"
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
