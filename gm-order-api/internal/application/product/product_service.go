package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ivanperez-dev/gm-order-api/internal/domain/product"
)

type ProductService struct {
	repo product.Repository
}

func NewProductService(repo product.Repository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, req CreateProductRequest) (*ProductResponse, error) {
	id := uuid.New().String()
	p, err := product.NewProduct(id, req.Name, req.Price, req.Stock)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, p); err != nil {
		return nil, errors.New("error saving product: " + err.Error())
	}

	return toResponse(p), nil
}

func (s *ProductService) GetByID(ctx context.Context, id string) (*ProductResponse, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("product not found: " + id)
	}
	return toResponse(p), nil
}

func (s *ProductService) ListAll(ctx context.Context, ids []string) ([]*ProductResponse, error) {
	var (
		products []*product.Product
		err      error
	)

	if len(ids) > 0 {
		products, err = s.repo.FindByIDs(ctx, ids)
	} else {
		products, err = s.repo.FindAll(ctx)
	}
	if err != nil {
		return nil, err
	}

	responses := make([]*ProductResponse, len(products))
	for i, p := range products {
		responses[i] = toResponse(p)
	}
	return responses, nil
}

func (s *ProductService) Update(ctx context.Context, id string, req UpdateProductRequest) (*ProductResponse, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("product not found: " + id)
	}

	if err := p.UpdateDetails(req.Name, req.Price, req.Stock); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, errors.New("error updating product: " + err.Error())
	}

	return toResponse(p), nil
}

func (s *ProductService) Delete(ctx context.Context, id string) error {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return errors.New("product not found: " + id)
	}

	return s.repo.Delete(ctx, id)
}
