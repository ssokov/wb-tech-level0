package kafka

import (
	"context"
	"encoding/json"
	"github.com/kimvlry/simple-order-service/internal/domain"
	"github.com/kimvlry/simple-order-service/internal/repo"
	"github.com/kimvlry/simple-order-service/internal/service/redis"
	"github.com/segmentio/kafka-go"
	"log"
)

type Consumer struct {
	repo        repo.OrderRepo
	redisClient *redis.Client
	reader      *kafka.Reader
}

func NewConsumer(repo repo.OrderRepo, redisClient *redis.Client, brokers []string, groupID, topic string) *Consumer {
	return &Consumer{
		repo:        repo,
		redisClient: redisClient,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: groupID,
			Topic:   topic,
		})}
}

func (c *Consumer) Consume() {
	defer c.reader.Close()
	ctx := context.Background()

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message from Kafka: %v", err)
			continue
		}
		log.Printf("Read message at topic/partition/offset %v/%v/%v: %v", m.Topic, m.Partition, m.Offset, string(m.Value))

		var order domain.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			log.Printf("Error unmarshalling message from Kafka: %v", err)
			continue
		}

		if err := order.ValidateOrder(); err != nil {
			log.Printf("Error validating order from Kafka: %v", err)
			continue
		}

		if err := c.repo.Save(ctx, &order); err != nil {
			log.Printf("Error saving order to DB: %v", err)
			continue
		}
		log.Printf("Order from message %s saved to DB with ID %v", string(m.Key), order.OrderUid)

		if err := c.redisClient.SaveOrder(ctx, &order); err != nil {
			log.Printf("Error caching order to Redis: %v", err)
			continue
		}
		log.Printf("Order %v from message %s cached to Redis", order.OrderUid, string(m.Key))
	}
}
