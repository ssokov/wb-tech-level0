package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/kimvlry/simple-order-service/internal/domain"
	"github.com/kimvlry/simple-order-service/internal/repo"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type Client struct {
	client        *redis.Client
	cacheDuration time.Duration
	orderRepo     repo.OrderRepo
}

func NewRedisClient(address string, duration time.Duration, repo repo.OrderRepo) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})
	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("couldn't connect to Redis: %v", err)
	}
	return &Client{
		client:        rdb,
		cacheDuration: duration,
		orderRepo:     repo,
	}
}

func (r *Client) SaveOrder(ctx context.Context, order *domain.Order) error {
	key := order.OrderUid

	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, r.cacheDuration).Err()
}

func (r *Client) GetOrder(ctx context.Context, orderUid string) (*domain.Order, error) {
	data, err := r.client.Get(ctx, orderUid).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var order domain.Order
	if err := json.Unmarshal([]byte(data), &order); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *Client) RestoreCache(ctx context.Context) error {
	orders, err := r.orderRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		if err := r.SaveOrder(ctx, &order); err != nil {
			return err
		}
	}
	return nil
}
