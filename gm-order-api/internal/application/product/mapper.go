package application

import "github.com/ivanperez-dev/gm-order-api/internal/domain/product"

func toResponse(p *product.Product) *ProductResponse {
	return &ProductResponse{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
