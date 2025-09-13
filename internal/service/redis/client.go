package redis

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/kimvlry/simple-order-service/internal/config"
	"github.com/kimvlry/simple-order-service/internal/domain"
	"github.com/kimvlry/simple-order-service/internal/interfaces"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	client        *redis.Client
	cacheDuration time.Duration
	orderRepo     interfaces.OrderRepository
}

func NewRedisClient(cfg config.RedisConfig, repo interfaces.OrderRepository) interfaces.Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("couldn't connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")

	return &Client{
		client:        rdb,
		cacheDuration: cfg.TTL,
		orderRepo:     repo,
	}
}

func (r *Client) SaveOrder(ctx context.Context, order *domain.Order) error {
	key := order.OrderUid

	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	log.Printf("Saving order %s to cache", key)
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
		log.Printf("error getting orders to restore to cache: %s", err)
		return err
	}

	for _, order := range orders {
		if err := r.SaveOrder(ctx, &order); err != nil {
			log.Printf("failed to save order %s to cache", order.OrderUid)
			return err
		}
	}
	log.Print("Redis cache restored")
	return nil
}

func (r *Client) Close() error {
	return r.client.Close()
}
