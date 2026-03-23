package product

import (
	"errors"
	"time"
)

type Product struct {
	ID        string
	Name      string
	Price     float64
	Stock     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProduct(id, name string, price float64, stock int) (*Product, error) {
	if id == "" {
		return nil, errors.New("productId is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}
	if price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}
	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	now := time.Now()
	return &Product{
		ID:        id,
		Name:      name,
		Price:     price,
		Stock:     stock,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (p *Product) UpdateDetails(name string, price float64, stock int) error {
	if name == "" {
		return errors.New("name is required")
	}
	if price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if stock < 0 {
		return errors.New("stock cannot be negative")
	}
	p.Name = name
	p.Price = price
	p.Stock = stock
	p.UpdatedAt = time.Now()
	return nil
}
