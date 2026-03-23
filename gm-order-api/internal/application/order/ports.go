package application

import "context"

type CreateOrderUseCase interface {
	Create(ctx context.Context, req CreateOrderRequest) (*OrderResponse, error)
}

type GetOrderUseCase interface {
	GetByID(ctx context.Context, id string) (*OrderResponse, error)
}

type ListOrdersUseCase interface {
	ListAll(ctx context.Context) ([]*OrderResponse, error)
}

var _ CreateOrderUseCase = (*OrderService)(nil)
var _ GetOrderUseCase = (*OrderService)(nil)
var _ ListOrdersUseCase = (*OrderService)(nil)
