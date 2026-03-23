package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ivanperez-dev/gm-order-api/internal/domain/customer"
)

type CustomerService struct {
	repo customer.Repository
}

func NewCustomerService(repo customer.Repository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) Create(ctx context.Context, req CreateCustomerRequest) (*CustomerResponse, error) {
	id := uuid.New().String()
	c, err := customer.NewCustomer(id, req.Name, req.Email)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, c); err != nil {
		return nil, errors.New("error saving customer: " + err.Error())
	}

	return toResponse(c), nil
}

func (s *CustomerService) GetByID(ctx context.Context, id string) (*CustomerResponse, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("customer not found: " + id)
	}
	return toResponse(c), nil
}

func (s *CustomerService) ListAll(ctx context.Context) ([]*CustomerResponse, error) {
	customers, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*CustomerResponse, len(customers))
	for i, c := range customers {
		responses[i] = toResponse(c)
	}
	return responses, nil
}

func (s *CustomerService) Update(ctx context.Context, id string, req UpdateCustomerRequest) (*CustomerResponse, error) {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("customer not found: " + id)
	}

	if err := c.UpdateDetails(req.Name, req.Email); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, errors.New("error updating customer: " + err.Error())
	}

	return toResponse(c), nil
}

func (s *CustomerService) Delete(ctx context.Context, id string) error {
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if c == nil {
		return errors.New("customer not found: " + id)
	}

	return s.repo.Delete(ctx, id)
}
