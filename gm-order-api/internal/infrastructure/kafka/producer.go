package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	domainorder "github.com/ivanperez-dev/gm-order-api/internal/domain/order"
	"github.com/segmentio/kafka-go"
)

type OrderCreatedEvent struct {
	OrderID    string           `json:"orderId"`
	CustomerID string           `json:"customerId"`
	Items      []OrderItemEvent `json:"items"`
	OccurredAt time.Time        `json:"occurredAt"`
}

type OrderItemEvent struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type KafkaEventProducer struct {
	writer *kafka.Writer
}

func NewKafkaEventProducer(writer *kafka.Writer) *KafkaEventProducer {
	return &KafkaEventProducer{writer: writer}
}

func (p *KafkaEventProducer) PublishOrderCreated(ctx context.Context, o *domainorder.Order) error {
	event := toEvent(o)

	payload, err := json.Marshal(event)
	if err != nil {
		return errors.New("error serializing order event: " + err.Error())
	}

	msg := kafka.Message{
		Key:   []byte(o.ID),
		Value: payload,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return errors.New("error publishing order event: " + err.Error())
	}

	log.Printf("[KAFKA] event published for order: %s", o.ID)
	return nil
}

func (p *KafkaEventProducer) Close() error {
	return p.writer.Close()
}

func toEvent(o *domainorder.Order) OrderCreatedEvent {
	items := make([]OrderItemEvent, len(o.Items))
	for i, item := range o.Items {
		items[i] = OrderItemEvent{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}
	return OrderCreatedEvent{
		OrderID:    o.ID,
		CustomerID: o.CustomerID,
		Items:      items,
		OccurredAt: time.Now(),
	}
}
