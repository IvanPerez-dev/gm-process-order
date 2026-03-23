package customer_test

import (
	"testing"

	"github.com/ivanperez-dev/gm-order-api/internal/domain/customer"
)

func TestNewCustomer_ShouldCreateCustomer_WhenDataIsValid(t *testing.T) {
	c, err := customer.NewCustomer("cust-001", "John Doe", "john@example.com")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if c.ID != "cust-001" {
		t.Errorf("expected ID cust-001, got %s", c.ID)
	}
	if c.Name != "John Doe" {
		t.Errorf("expected name John Doe, got %s", c.Name)
	}
	if c.Email != "john@example.com" {
		t.Errorf("expected email john@example.com, got %s", c.Email)
	}
	if !c.IsActive {
		t.Error("expected customer to be active")
	}
}

func TestNewCustomer_ShouldReturnError_WhenIDIsEmpty(t *testing.T) {
	c, err := customer.NewCustomer("", "John Doe", "john@example.com")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if c != nil {
		t.Error("expected nil customer")
	}
}

func TestNewCustomer_ShouldReturnError_WhenNameIsEmpty(t *testing.T) {
	c, err := customer.NewCustomer("cust-001", "", "john@example.com")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if c != nil {
		t.Error("expected nil customer")
	}
}

func TestNewCustomer_ShouldReturnError_WhenEmailIsInvalid(t *testing.T) {
	invalidEmails := []string{"notanemail", "missing-at.com", "@nodomain", ""}

	for _, email := range invalidEmails {
		c, err := customer.NewCustomer("cust-001", "John Doe", email)

		if err == nil {
			t.Errorf("expected error for email %q, got nil", email)
		}
		if c != nil {
			t.Errorf("expected nil customer for email %q", email)
		}
	}
}

func TestUpdateDetails_ShouldUpdateCustomer_WhenDataIsValid(t *testing.T) {
	c, _ := customer.NewCustomer("cust-001", "John Doe", "john@example.com")

	err := c.UpdateDetails("Jane Doe", "jane@example.com")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if c.Name != "Jane Doe" {
		t.Errorf("expected name Jane Doe, got %s", c.Name)
	}
	if c.Email != "jane@example.com" {
		t.Errorf("expected email jane@example.com, got %s", c.Email)
	}
}

func TestUpdateDetails_ShouldReturnError_WhenNameIsEmpty(t *testing.T) {
	c, _ := customer.NewCustomer("cust-001", "John Doe", "john@example.com")

	err := c.UpdateDetails("", "john@example.com")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateDetails_ShouldReturnError_WhenEmailIsInvalid(t *testing.T) {
	c, _ := customer.NewCustomer("cust-001", "John Doe", "john@example.com")

	err := c.UpdateDetails("John Doe", "invalid-email")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
