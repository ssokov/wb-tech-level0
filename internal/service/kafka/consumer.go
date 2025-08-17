package kafka

import (
	"context"
	"encoding/json"
	kafka "github.com/segmentio/kafka-go"
	"simple-order-service/internal/domain"
	"simple-order-service/internal/repo"
	"simple-order-service/internal/service/redis"

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
			Brokers:        brokers,
			GroupID:        groupID,
			Topic:          topic,
			StartOffset:    kafka.FirstOffset,
			CommitInterval: 0,
		}),
	}
}

func (c *Consumer) Consume() {
	defer c.reader.Close()
	ctx := context.Background()

	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("Error fetching message from Kafka: %v", err)
			continue
		}

		log.Printf("Read message at topic/partition/offset %v/%v/%v: %v", m.Topic, m.Partition, m.Offset, string(m.Value))

		var order domain.Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			log.Printf("Error unmarshalling message from Kafka: %v", err)
			if err := c.reader.CommitMessages(ctx, m); err != nil {
				log.Printf("Failed to commit offset after unmarshalling error: %v", err)
			}
			continue
		}

		if err := order.ValidateOrder(); err != nil {
			log.Printf("Error validating order from Kafka: %v", err)
			if err := c.reader.CommitMessages(ctx, m); err != nil {
				log.Printf("Failed to commit offset after validation error: %v", err)
			}
			continue
		}

		if err := c.repo.Save(ctx, &order); err != nil {
			log.Printf("Error saving order to DB: %v", err)
			continue
		}

		if err := c.redisClient.SaveOrder(ctx, &order); err != nil {
			log.Printf("Error caching order to Redis: %v", err)
			continue
		}

		if err := c.reader.CommitMessages(ctx, m); err != nil {
			log.Printf("Failed to commit offset: %v", err)
		} else {
			log.Printf("Order from message %s saved and cached, offset committed", string(m.Key))
		}
	}
}
