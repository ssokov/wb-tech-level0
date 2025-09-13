package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/kimvlry/simple-order-service/internal/config"
	"github.com/kimvlry/simple-order-service/internal/domain"
	"github.com/kimvlry/simple-order-service/internal/interfaces"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	repo        interfaces.OrderRepository
	redisClient interfaces.Cache
	reader      *kafka.Reader
	dlqWriter   *kafka.Writer
	config      config.KafkaConfig
}

func NewConsumer(repo interfaces.OrderRepository, redisClient interfaces.Cache, cfg config.KafkaConfig) interfaces.MessageConsumer {
	return &Consumer{
		repo:        repo,
		redisClient: redisClient,
		config:      cfg,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        cfg.Brokers,
			GroupID:        cfg.GroupID,
			Topic:          cfg.Topic,
			StartOffset:    kafka.FirstOffset,
			CommitInterval: 0,
		}),
		dlqWriter: &kafka.Writer{
			Addr:     kafka.TCP(cfg.Brokers...),
			Topic:    cfg.DLQTopic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	defer c.reader.Close()
	defer c.dlqWriter.Close()

	retryCount := 0
	maxRetries := 3

	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer context cancelled")
			return ctx.Err()
		default:
		}

		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("Error fetching message from Kafka: %v", err)
			continue
		}

		log.Printf("Read message at topic/partition/offset %v/%v/%v: %v", m.Topic, m.Partition, m.Offset, string(m.Value))

		// Process message with retry logic
		if err := c.processMessage(ctx, m, retryCount, maxRetries); err != nil {
			log.Printf("Failed to process message after %d retries: %v", retryCount, err)
			// Send to DLQ
			if err := c.sendToDLQ(ctx, m, err); err != nil {
				log.Printf("Failed to send message to DLQ: %v", err)
			}
		}

		// Commit message
		if err := c.reader.CommitMessages(ctx, m); err != nil {
			log.Printf("Failed to commit offset: %v", err)
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, m kafka.Message, retryCount, maxRetries int) error {
	var order domain.Order
	if err := json.Unmarshal(m.Value, &order); err != nil {
		log.Printf("Error unmarshalling message from Kafka: %v", err)
		return err
	}

	if err := order.ValidateOrder(); err != nil {
		log.Printf("Error validating order from Kafka: %v", err)
		return err
	}

	if err := c.repo.Save(ctx, &order); err != nil {
		log.Printf("Error saving order to DB: %v", err)
		if retryCount < maxRetries {
			time.Sleep(time.Duration(retryCount+1) * time.Second)
			return c.processMessage(ctx, m, retryCount+1, maxRetries)
		}
		return err
	}

	if err := c.redisClient.SaveOrder(ctx, &order); err != nil {
		log.Printf("Error caching order to Redis: %v", err)
		// Don't retry cache errors, just log them
	}

	log.Printf("Order from message %s saved and cached", string(m.Key))
	return nil
}

func (c *Consumer) sendToDLQ(ctx context.Context, m kafka.Message, err error) error {
	dlqMessage := kafka.Message{
		Key:   m.Key,
		Value: m.Value,
		Headers: []kafka.Header{
			{Key: "original-topic", Value: []byte(m.Topic)},
			{Key: "error", Value: []byte(err.Error())},
			{Key: "timestamp", Value: []byte(time.Now().Format(time.RFC3339))},
		},
	}

	return c.dlqWriter.WriteMessages(ctx, dlqMessage)
}

func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		return err
	}
	return c.dlqWriter.Close()
}
