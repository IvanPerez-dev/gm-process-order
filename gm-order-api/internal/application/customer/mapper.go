package application

import "github.com/ivanperez-dev/gm-order-api/internal/domain/customer"

func toResponse(c *customer.Customer) *CustomerResponse {
	return &CustomerResponse{
		ID:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		IsActive:  c.IsActive,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
