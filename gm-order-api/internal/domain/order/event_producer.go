package order

import "context"

type EventProducer interface {
	PublishOrderCreated(ctx context.Context, order *Order) error
}
