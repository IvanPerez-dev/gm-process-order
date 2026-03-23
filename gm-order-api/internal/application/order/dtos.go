package application

import "time"

type OrderItemRequest struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int    `json:"quantity"  binding:"required,min=1"`
}

type CreateOrderRequest struct {
	CustomerID string             `json:"customerId" binding:"required"`
	Items      []OrderItemRequest `json:"items"      binding:"required,min=1"`
}

type OrderItemResponse struct {
	ProductID string  `json:"productId"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

type OrderCustomerResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type OrderResponse struct {
	ID          string              `json:"id"`
	CustomerID  string              `json:"customerId"`
	Items       []OrderItemResponse `json:"items"`
	Status      string              `json:"status"`
	RetryCount  int                 `json:"retryCount"`
	CreatedAt   time.Time           `json:"createdAt"`
	ProcessedAt *time.Time          `json:"processedAt,omitempty"`
}
