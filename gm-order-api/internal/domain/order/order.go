package order

import (
	"errors"
	"time"
)

type OrderStatus string

const (
	StatusPending    OrderStatus = "PENDING"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusProcessed  OrderStatus = "PROCESSED"
	StatusFailed     OrderStatus = "FAILED"
)

type OrderItem struct {
	ProductID string
	Name      string
	Price     float64
	Quantity  int
}

func (i OrderItem) Validate() error {
	if i.ProductID == "" {
		return errors.New("productId is required")
	}
	if i.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	return nil
}

type Order struct {
	ID          string
	CustomerID  string
	Items       []OrderItem
	Status      OrderStatus
	RetryCount  int
	ErrorMsg    string
	CreatedAt   time.Time
	ProcessedAt *time.Time
}

func NewOrder(id, customerID string, items []OrderItem) (*Order, error) {
	if id == "" {
		return nil, errors.New("orderId is required")
	}
	if customerID == "" {
		return nil, errors.New("customerId is required")
	}
	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	for _, item := range items {
		if err := item.Validate(); err != nil {
			return nil, err
		}
	}

	return &Order{
		ID:         id,
		CustomerID: customerID,
		Items:      items,
		Status:     StatusPending,
		RetryCount: 0,
		CreatedAt:  time.Now(),
	}, nil
}

func (o *Order) MarkAsProcessing() error {
	if o.Status != StatusPending {
		return errors.New("only PENDING orders can be marked as PROCESSING")
	}
	o.Status = StatusProcessing
	return nil
}

func (o *Order) MarkAsProcessed(enrichedItems []OrderItem) error {
	if o.Status != StatusProcessing {
		return errors.New("only PROCESSING orders can be marked as PROCESSED")
	}
	now := time.Now()
	o.Items = enrichedItems
	o.Status = StatusProcessed
	o.ProcessedAt = &now
	return nil
}

func (o *Order) MarkAsFailed(reason string) {
	o.Status = StatusFailed
	o.ErrorMsg = reason
}

func (o *Order) IncrementRetry() {
	o.RetryCount++
}

func (o *Order) HasExceededRetries(maxRetries int) bool {
	return o.RetryCount >= maxRetries
}
