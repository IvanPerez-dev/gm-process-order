package application

import "context"

type CreateCustomerUseCase interface {
	Create(ctx context.Context, req CreateCustomerRequest) (*CustomerResponse, error)
}

type GetCustomerUseCase interface {
	GetByID(ctx context.Context, id string) (*CustomerResponse, error)
}

type ListCustomersUseCase interface {
	ListAll(ctx context.Context) ([]*CustomerResponse, error)
}

type UpdateCustomerUseCase interface {
	Update(ctx context.Context, id string, req UpdateCustomerRequest) (*CustomerResponse, error)
}

type DeleteCustomerUseCase interface {
	Delete(ctx context.Context, id string) error
}

var _ CreateCustomerUseCase = (*CustomerService)(nil)
var _ GetCustomerUseCase = (*CustomerService)(nil)
var _ ListCustomersUseCase = (*CustomerService)(nil)
var _ UpdateCustomerUseCase = (*CustomerService)(nil)
var _ DeleteCustomerUseCase = (*CustomerService)(nil)
