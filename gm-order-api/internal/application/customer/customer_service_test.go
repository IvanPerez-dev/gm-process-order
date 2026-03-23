package application

import (
	"context"
	"errors"
	"testing"

	"github.com/ivanperez-dev/gm-order-api/internal/domain/customer"
)

// --- mock ---

type mockCustomerRepo struct {
	saveFn     func(ctx context.Context, c *customer.Customer) error
	findByIDFn func(ctx context.Context, id string) (*customer.Customer, error)
	findAllFn  func(ctx context.Context) ([]*customer.Customer, error)
	updateFn   func(ctx context.Context, c *customer.Customer) error
	deleteFn   func(ctx context.Context, id string) error
}

func (m *mockCustomerRepo) Save(ctx context.Context, c *customer.Customer) error {
	return m.saveFn(ctx, c)
}

func (m *mockCustomerRepo) FindByID(ctx context.Context, id string) (*customer.Customer, error) {
	return m.findByIDFn(ctx, id)
}

func (m *mockCustomerRepo) FindAll(ctx context.Context) ([]*customer.Customer, error) {
	return m.findAllFn(ctx)
}

func (m *mockCustomerRepo) Update(ctx context.Context, c *customer.Customer) error {
	return m.updateFn(ctx, c)
}

func (m *mockCustomerRepo) Delete(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}

// helpers

func newValidCustomer(t *testing.T) *customer.Customer {
	t.Helper()
	c, err := customer.NewCustomer("cust-001", "John Doe", "john@example.com")
	if err != nil {
		t.Fatalf("failed to build test customer: %v", err)
	}
	return c
}

// --- Create ---

func TestCustomerCreate_ShouldReturnCustomer_WhenDataIsValid(t *testing.T) {
	repo := &mockCustomerRepo{
		saveFn: func(_ context.Context, _ *customer.Customer) error { return nil },
	}
	svc := NewCustomerService(repo)

	resp, err := svc.Create(context.Background(), CreateCustomerRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Name != "John Doe" {
		t.Errorf("expected name John Doe, got %s", resp.Name)
	}
	if resp.Email != "john@example.com" {
		t.Errorf("expected email john@example.com, got %s", resp.Email)
	}
	if !resp.IsActive {
		t.Error("expected customer to be active")
	}
	if resp.ID == "" {
		t.Error("expected non-empty ID")
	}
}

func TestCustomerCreate_ShouldReturnError_WhenDomainValidationFails(t *testing.T) {
	repo := &mockCustomerRepo{}
	svc := NewCustomerService(repo)

	_, err := svc.Create(context.Background(), CreateCustomerRequest{
		Name:  "",
		Email: "john@example.com",
	})

	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestCustomerCreate_ShouldReturnError_WhenEmailIsInvalid(t *testing.T) {
	repo := &mockCustomerRepo{}
	svc := NewCustomerService(repo)

	_, err := svc.Create(context.Background(), CreateCustomerRequest{
		Name:  "John Doe",
		Email: "not-an-email",
	})

	if err == nil {
		t.Fatal("expected error for invalid email, got nil")
	}
}

func TestCustomerCreate_ShouldReturnError_WhenRepoSaveFails(t *testing.T) {
	repo := &mockCustomerRepo{
		saveFn: func(_ context.Context, _ *customer.Customer) error {
			return errors.New("db connection error")
		},
	}
	svc := NewCustomerService(repo)

	_, err := svc.Create(context.Background(), CreateCustomerRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	})

	if err == nil {
		t.Fatal("expected error when repo save fails, got nil")
	}
}

// --- GetByID ---

func TestCustomerGetByID_ShouldReturnCustomer_WhenFound(t *testing.T) {
	c := newValidCustomer(t)
	repo := &mockCustomerRepo{
		findByIDFn: func(_ context.Context, _ string) (*customer.Customer, error) { return c, nil },
	}
	svc := NewCustomerService(repo)

	resp, err := svc.GetByID(context.Background(), "cust-001")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.ID != "cust-001" {
		t.Errorf("expected ID cust-001, got %s", resp.ID)
	}
}

func TestCustomerGetByID_ShouldReturnError_WhenNotFound(t *testing.T) {
	repo := &mockCustomerRepo{
		findByIDFn: func(_ context.Context, _ string) (*customer.Customer, error) { return nil, nil },
	}
	svc := NewCustomerService(repo)

	_, err := svc.GetByID(context.Background(), "cust-999")

	if err == nil {
		t.Fatal("expected error when customer not found, got nil")
	}
}

func TestCustomerGetByID_ShouldReturnError_WhenRepoFails(t *testing.T) {
	repo := &mockCustomerRepo{
		findByIDFn: func(_ context.Context, _ string) (*customer.Customer, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewCustomerService(repo)

	_, err := svc.GetByID(context.Background(), "cust-001")

	if err == nil {
		t.Fatal("expected error when repo fails, got nil")
	}
}

// --- ListAll ---

func TestCustomerListAll_ShouldReturnAllCustomers(t *testing.T) {
	c1 := newValidCustomer(t)
	c2, _ := customer.NewCustomer("cust-002", "Jane Doe", "jane@example.com")
	repo := &mockCustomerRepo{
		findAllFn: func(_ context.Context) ([]*customer.Customer, error) {
			return []*customer.Customer{c1, c2}, nil
		},
	}
	svc := NewCustomerService(repo)

	resp, err := svc.ListAll(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(resp) != 2 {
		t.Errorf("expected 2 customers, got %d", len(resp))
	}
}

func TestCustomerListAll_ShouldReturnEmptySlice_WhenNoCustomers(t *testing.T) {
	repo := &mockCustomerRepo{
		findAllFn: func(_ context.Context) ([]*customer.Customer, error) {
			return []*customer.Customer{}, nil
		},
	}
	svc := NewCustomerService(repo)

	resp, err := svc.ListAll(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(resp) != 0 {
		t.Errorf("expected empty slice, got %d items", len(resp))
	}
}

func TestCustomerListAll_ShouldReturnError_WhenRepoFails(t *testing.T) {
	repo := &mockCustomerRepo{
		findAllFn: func(_ context.Context) ([]*customer.Customer, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewCustomerService(repo)

	_, err := svc.ListAll(context.Background())

	if err == nil {
		t.Fatal("expected error when repo fails, got nil")
	}
}

// --- Update ---

func TestCustomerUpdate_ShouldReturnUpdatedCustomer_WhenValid(t *testing.T) {
	c := newValidCustomer(t)
	repo := &mockCustomerRepo{
		findByIDFn: func(_ context.Context, _ string) (*customer.Customer, error) { return c, nil },
		updateFn:   func(_ context.Context, _ *customer.Customer) error { return nil },
	}
	svc := NewCustomerService(repo)

	resp, err := svc.Update(context.Background(), "cust-001", UpdateCustomerRequest{
		Name:  "Jane Doe",
		Email: "jane@example.com",
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.Name != "Jane Doe" {
		t.Errorf("expected name Jane Doe, got %s", resp.Name)
	}
	if resp.Email != "jane@example.com" {
		t.Errorf("expected email jane@example.com, got %s", resp.Email)
	}
}

func TestCustomerUpdate_ShouldReturnError_WhenCustomerNotFound(t *testing.T) {
	repo := &mockCustomerRepo{
		findByIDFn: func(_ context.Context, _ string) (*customer.Customer, error) { return nil, nil },
	}
	svc := NewCustomerService(repo)

	_, err := svc.Update(context.Background(), "cust-999", UpdateCustomerRequest{
		Name:  "Jane Doe",
		Email: "jane@example.com",
	})

	if err == nil {
		t.Fatal("expected error when customer not found, got nil")
	}
}

func TestCustomerUpdate_ShouldReturnError_WhenDomainValidationFails(t *testing.T) {
	c := newValidCustomer(t)
	repo := &mockCustomerRepo{
		findByIDFn: func(_ context.Context, _ string) (*customer.Customer, error) { return c, nil },
	}
	svc := NewCustomerService(repo)

	_, err := svc.Update(context.Background(), "cust-001", UpdateCustomerRequest{
		Name:  "",
		Email: "jane@example.com",
	})

	if err == nil {
		t.Fatal("expected error for invalid name, got nil")
	}
}

func TestCustomerUpdate_ShouldReturnError_WhenRepoUpdateFails(t *testing.T) {
	c := newValidCustomer(t)
	repo := &mockCustomerRepo{
		findByIDFn: func(_ context.Context, _ string) (*customer.Customer, error) { return c, nil },
		updateFn: func(_ context.Context, _ *customer.Customer) error {
			return errors.New("db error")
		},
	}
	svc := NewCustomerService(repo)

	_, err := svc.Update(context.Background(), "cust-001", UpdateCustomerRequest{
		Name:  "Jane Doe",
		Email: "jane@example.com",
	})

	if err == nil {
		t.Fatal("expected error when repo update fails, got nil")
	}
}

// --- Delete ---

func TestCustomerDelete_ShouldSucceed_WhenCustomerExists(t *testing.T) {
	c := newValidCustomer(t)
	repo := &mockCustomerRepo{
		findByIDFn: func(_ context.Context, _ string) (*customer.Customer, error) { return c, nil },
		deleteFn:   func(_ context.Context, _ string) error { return nil },
	}
	svc := NewCustomerService(repo)

	err := svc.Delete(context.Background(), "cust-001")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestCustomerDelete_ShouldReturnError_WhenCustomerNotFound(t *testing.T) {
	repo := &mockCustomerRepo{
		findByIDFn: func(_ context.Context, _ string) (*customer.Customer, error) { return nil, nil },
	}
	svc := NewCustomerService(repo)

	err := svc.Delete(context.Background(), "cust-999")

	if err == nil {
		t.Fatal("expected error when customer not found, got nil")
	}
}

func TestCustomerDelete_ShouldReturnError_WhenRepoFindFails(t *testing.T) {
	repo := &mockCustomerRepo{
		findByIDFn: func(_ context.Context, _ string) (*customer.Customer, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewCustomerService(repo)

	err := svc.Delete(context.Background(), "cust-001")

	if err == nil {
		t.Fatal("expected error when repo find fails, got nil")
	}
}
