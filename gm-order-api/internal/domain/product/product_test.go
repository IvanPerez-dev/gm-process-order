package product_test

import (
	"testing"

	"github.com/ivanperez-dev/gm-order-api/internal/domain/product"
)

func TestNewProduct_ShouldCreateProduct_WhenDataIsValid(t *testing.T) {
	p, err := product.NewProduct("prod-001", "Laptop", 999.99, 10)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if p.ID != "prod-001" {
		t.Errorf("expected ID prod-001, got %s", p.ID)
	}
	if p.Name != "Laptop" {
		t.Errorf("expected name Laptop, got %s", p.Name)
	}
	if p.Price != 999.99 {
		t.Errorf("expected price 999.99, got %f", p.Price)
	}
	if p.Stock != 10 {
		t.Errorf("expected stock 10, got %d", p.Stock)
	}
}

func TestNewProduct_ShouldCreateProduct_WhenStockIsZero(t *testing.T) {
	p, err := product.NewProduct("prod-001", "Laptop", 999.99, 0)

	if err != nil {
		t.Fatalf("expected no error for zero stock, got: %v", err)
	}
	if p.Stock != 0 {
		t.Errorf("expected stock 0, got %d", p.Stock)
	}
}

func TestNewProduct_ShouldReturnError_WhenIDIsEmpty(t *testing.T) {
	p, err := product.NewProduct("", "Laptop", 999.99, 10)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if p != nil {
		t.Error("expected nil product")
	}
}

func TestNewProduct_ShouldReturnError_WhenNameIsEmpty(t *testing.T) {
	p, err := product.NewProduct("prod-001", "", 999.99, 10)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if p != nil {
		t.Error("expected nil product")
	}
}

func TestNewProduct_ShouldReturnError_WhenPriceIsZero(t *testing.T) {
	p, err := product.NewProduct("prod-001", "Laptop", 0, 10)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if p != nil {
		t.Error("expected nil product")
	}
}

func TestNewProduct_ShouldReturnError_WhenPriceIsNegative(t *testing.T) {
	p, err := product.NewProduct("prod-001", "Laptop", -1.0, 10)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if p != nil {
		t.Error("expected nil product")
	}
}

func TestNewProduct_ShouldReturnError_WhenStockIsNegative(t *testing.T) {
	p, err := product.NewProduct("prod-001", "Laptop", 999.99, -1)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if p != nil {
		t.Error("expected nil product")
	}
}

func TestUpdateDetails_ShouldUpdateProduct_WhenDataIsValid(t *testing.T) {
	p, _ := product.NewProduct("prod-001", "Laptop", 999.99, 10)

	err := p.UpdateDetails("Gaming Laptop", 1299.99, 5)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if p.Name != "Gaming Laptop" {
		t.Errorf("expected name Gaming Laptop, got %s", p.Name)
	}
	if p.Price != 1299.99 {
		t.Errorf("expected price 1299.99, got %f", p.Price)
	}
	if p.Stock != 5 {
		t.Errorf("expected stock 5, got %d", p.Stock)
	}
}

func TestUpdateDetails_ShouldReturnError_WhenNameIsEmpty(t *testing.T) {
	p, _ := product.NewProduct("prod-001", "Laptop", 999.99, 10)

	err := p.UpdateDetails("", 999.99, 10)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateDetails_ShouldReturnError_WhenPriceIsZero(t *testing.T) {
	p, _ := product.NewProduct("prod-001", "Laptop", 999.99, 10)

	err := p.UpdateDetails("Laptop", 0, 10)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateDetails_ShouldReturnError_WhenPriceIsNegative(t *testing.T) {
	p, _ := product.NewProduct("prod-001", "Laptop", 999.99, 10)

	err := p.UpdateDetails("Laptop", -5.0, 10)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateDetails_ShouldReturnError_WhenStockIsNegative(t *testing.T) {
	p, _ := product.NewProduct("prod-001", "Laptop", 999.99, 10)

	err := p.UpdateDetails("Laptop", 999.99, -1)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
