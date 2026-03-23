package application

import (
	"github.com/google/uuid"
	"github.com/ivanperez-dev/gm-order-api/internal/domain/order"
)

func generateID() string {
	return uuid.New().String()
}

func toResponse(o *order.Order) *OrderResponse {
	items := make([]OrderItemResponse, len(o.Items))
	for i, item := range o.Items {
		items[i] = OrderItemResponse{
			ProductID: item.ProductID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
			Total:     float64(item.Quantity) * item.Price,
		}
	}
	return &OrderResponse{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Items:       items,
		Status:      string(o.Status),
		CreatedAt:   o.CreatedAt,
		ProcessedAt: o.ProcessedAt,
	}
}
