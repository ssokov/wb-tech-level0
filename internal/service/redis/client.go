package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/kimvlry/simple-order-service/internal/domain"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type RedisClient struct {
	client        *redis.Client
	ctx           context.Context
	cacheDuration time.Duration
}

func NewRedisClient(address string, duration time.Duration) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})
	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("couldn't connect to Redis: %v", err)
	}
	return &RedisClient{
		client:        rdb,
		ctx:           ctx,
		cacheDuration: duration,
	}
}

func (r *RedisClient) SaveOrder(ctx context.Context, order *domain.Order) error {
	key := order.OrderUid

	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, r.cacheDuration).Err()
}

func (r *RedisClient) GetOrder(ctx context.Context, orderUid string) (*domain.Order, error) {
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
