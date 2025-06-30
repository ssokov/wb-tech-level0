package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kimvlry/simple-order-service/internal/domain"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer kafka.Writer
}

func NewProducer(address, topic string) *Producer {
	return &Producer{
		writer: kafka.Writer{
			Addr:     kafka.TCP(address),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) PublishOrder(ctx context.Context, order domain.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(order.OrderUid),
		Value: data,
		Time:  time.Now(),
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
