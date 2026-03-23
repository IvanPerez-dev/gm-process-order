package order_test

import (
	"testing"

	"github.com/ivanperez-dev/gm-order-api/internal/domain/order"
)

func TestNewOrder_ShouldCreateOrder_WhenDataIsValid(t *testing.T) {

	items := []order.OrderItem{
		{ProductID: "product-001", Quantity: 2},
	}

	o, err := order.NewOrder("order-123", "customer-456", items)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if o.ID != "order-123" {
		t.Errorf("expected ID order-123, got %s", o.ID)
	}
	if o.Status != order.StatusPending {
		t.Errorf("expected status PENDING, got %s", o.Status)
	}
	if len(o.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(o.Items))
	}
}

func TestNewOrder_ShouldReturnError_WhenCustomerIDIsEmpty(t *testing.T) {
	// Arrange
	items := []order.OrderItem{
		{ProductID: "product-001", Quantity: 1},
	}

	o, err := order.NewOrder("order-123", "", items)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if o != nil {
		t.Error("expected nil order")
	}
}

func TestNewOrder_ShouldReturnError_WhenItemsIsEmpty(t *testing.T) {
	// Act
	o, err := order.NewOrder("order-123", "customer-456", []order.OrderItem{})

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if o != nil {
		t.Error("expected nil order")
	}
}

func TestMarkAsProcessing_ShouldTransition_WhenStatusIsPending(t *testing.T) {
	// Arrange
	o, _ := order.NewOrder("order-123", "customer-456", []order.OrderItem{
		{ProductID: "product-001", Quantity: 1},
	})

	err := o.MarkAsProcessing()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if o.Status != order.StatusProcessing {
		t.Errorf("expected status PROCESSING, got %s", o.Status)
	}
}

func TestMarkAsProcessing_ShouldReturnError_WhenStatusIsNotPending(t *testing.T) {

	o, _ := order.NewOrder("order-123", "customer-456", []order.OrderItem{
		{ProductID: "product-001", Quantity: 1},
	})
	o.MarkAsProcessing()

	err := o.MarkAsProcessing()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMarkAsProcessed_ShouldTransition_WhenStatusIsProcessing(t *testing.T) {

	o, _ := order.NewOrder("order-123", "customer-456", []order.OrderItem{
		{ProductID: "product-001", Quantity: 1},
	})
	o.MarkAsProcessing()

	enrichedItems := []order.OrderItem{
		{ProductID: "product-001", Name: "Laptop", Price: 999, Quantity: 1},
	}

	err := o.MarkAsProcessed(enrichedItems)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if o.Status != order.StatusProcessed {
		t.Errorf("expected status PROCESSED, got %s", o.Status)
	}
	if o.Items[0].Name != "Laptop" {
		t.Errorf("expected item name Laptop, got %s", o.Items[0].Name)
	}
	if o.ProcessedAt == nil {
		t.Error("expected processedAt to be set")
	}
}

func TestHasExceededRetries_ShouldReturnTrue_WhenMaxRetriesReached(t *testing.T) {

	o, _ := order.NewOrder("order-123", "customer-456", []order.OrderItem{
		{ProductID: "product-001", Quantity: 1},
	})

	o.IncrementRetry()
	o.IncrementRetry()
	o.IncrementRetry()

	if !o.HasExceededRetries(3) {
		t.Error("expected HasExceededRetries to return true")
	}
}
