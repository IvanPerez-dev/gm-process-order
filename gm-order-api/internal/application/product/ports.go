package application

import "context"

type CreateProductUseCase interface {
	Create(ctx context.Context, req CreateProductRequest) (*ProductResponse, error)
}

type GetProductUseCase interface {
	GetByID(ctx context.Context, id string) (*ProductResponse, error)
}

type ListProductsUseCase interface {
	ListAll(ctx context.Context, ids []string) ([]*ProductResponse, error)
}

type UpdateProductUseCase interface {
	Update(ctx context.Context, id string, req UpdateProductRequest) (*ProductResponse, error)
}

type DeleteProductUseCase interface {
	Delete(ctx context.Context, id string) error
}

var _ CreateProductUseCase = (*ProductService)(nil)
var _ GetProductUseCase = (*ProductService)(nil)
var _ ListProductsUseCase = (*ProductService)(nil)
var _ UpdateProductUseCase = (*ProductService)(nil)
var _ DeleteProductUseCase = (*ProductService)(nil)
