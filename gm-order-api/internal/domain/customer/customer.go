package customer

import (
	"errors"
	"strings"
	"time"
)

type Customer struct {
	ID        string
	Name      string
	Email     string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCustomer(id, name, email string) (*Customer, error) {
	if id == "" {
		return nil, errors.New("customerId is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}
	if !isValidEmail(email) {
		return nil, errors.New("email is invalid")
	}

	now := time.Now()
	return &Customer{
		ID:        id,
		Name:      name,
		Email:     email,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (c *Customer) UpdateDetails(name, email string) error {
	if name == "" {
		return errors.New("name is required")
	}
	if !isValidEmail(email) {
		return errors.New("email is invalid")
	}
	c.Name = name
	c.Email = email
	c.UpdatedAt = time.Now()
	return nil
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
