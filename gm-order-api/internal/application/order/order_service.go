package application

import (
	"context"
	"errors"

	"github.com/ivanperez-dev/gm-order-api/internal/domain/order"
)

type OrderService struct {
	repo          order.Repository
	eventProducer order.EventProducer
}

func NewOrderService(repo order.Repository, eventProducer order.EventProducer) *OrderService {
	return &OrderService{
		repo:          repo,
		eventProducer: eventProducer,
	}
}

func (s *OrderService) Create(ctx context.Context, req CreateOrderRequest) (*OrderResponse, error) {
	items := make([]order.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = order.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}
	id := generateID()
	newOrder, err := order.NewOrder(id, req.CustomerID, items)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, newOrder); err != nil {
		return nil, errors.New("error saving order: " + err.Error())
	}

	if err := s.eventProducer.PublishOrderCreated(ctx, newOrder); err != nil {

		return nil, errors.New("order saved but failed to publish event: " + err.Error())
	}

	return toResponse(newOrder), nil
}

func (s *OrderService) GetByID(ctx context.Context, id string) (*OrderResponse, error) {
	o, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if o == nil {
		return nil, errors.New("order not found: " + id)
	}
	return toResponse(o), nil
}

func (s *OrderService) ListAll(ctx context.Context) ([]*OrderResponse, error) {
	orders, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*OrderResponse, len(orders))
	for i, o := range orders {
		responses[i] = toResponse(o)
	}
	return responses, nil
}
