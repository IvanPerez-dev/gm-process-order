package application

import (
	"context"
	"errors"
	"testing"

	"github.com/ivanperez-dev/gm-order-api/internal/domain/order"
)

// --- mocks ---

type mockOrderRepo struct {
	saveFn     func(ctx context.Context, o *order.Order) error
	findByIDFn func(ctx context.Context, id string) (*order.Order, error)
	findAllFn  func(ctx context.Context) ([]*order.Order, error)
	updateFn   func(ctx context.Context, o *order.Order) error
}

func (m *mockOrderRepo) Save(ctx context.Context, o *order.Order) error {
	return m.saveFn(ctx, o)
}

func (m *mockOrderRepo) FindByID(ctx context.Context, id string) (*order.Order, error) {
	return m.findByIDFn(ctx, id)
}

func (m *mockOrderRepo) FindAll(ctx context.Context) ([]*order.Order, error) {
	return m.findAllFn(ctx)
}

func (m *mockOrderRepo) Update(ctx context.Context, o *order.Order) error {
	return m.updateFn(ctx, o)
}

type mockEventProducer struct {
	publishFn func(ctx context.Context, o *order.Order) error
}

func (m *mockEventProducer) PublishOrderCreated(ctx context.Context, o *order.Order) error {
	return m.publishFn(ctx, o)
}

// helpers

func validCreateRequest() CreateOrderRequest {
	return CreateOrderRequest{
		CustomerID: "cust-001",
		Items: []OrderItemRequest{
			{ProductID: "prod-001", Quantity: 2},
		},
	}
}

func newValidOrder(t *testing.T) *order.Order {
	t.Helper()
	o, err := order.NewOrder("order-001", "cust-001", []order.OrderItem{
		{ProductID: "prod-001", Quantity: 2},
	})
	if err != nil {
		t.Fatalf("failed to build test order: %v", err)
	}
	return o
}

// --- Create ---

func TestOrderCreate_ShouldReturnOrder_WhenDataIsValid(t *testing.T) {
	repo := &mockOrderRepo{
		saveFn: func(_ context.Context, _ *order.Order) error { return nil },
	}
	producer := &mockEventProducer{
		publishFn: func(_ context.Context, _ *order.Order) error { return nil },
	}
	svc := NewOrderService(repo, producer)

	resp, err := svc.Create(context.Background(), validCreateRequest())

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.CustomerID != "cust-001" {
		t.Errorf("expected customerID cust-001, got %s", resp.CustomerID)
	}
	if resp.ID == "" {
		t.Error("expected non-empty ID")
	}
	if resp.Status != string(order.StatusPending) {
		t.Errorf("expected status PENDING, got %s", resp.Status)
	}
	if len(resp.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(resp.Items))
	}
}

func TestOrderCreate_ShouldReturnError_WhenCustomerIDIsEmpty(t *testing.T) {
	repo := &mockOrderRepo{}
	producer := &mockEventProducer{}
	svc := NewOrderService(repo, producer)

	_, err := svc.Create(context.Background(), CreateOrderRequest{
		CustomerID: "",
		Items:      []OrderItemRequest{{ProductID: "prod-001", Quantity: 1}},
	})

	if err == nil {
		t.Fatal("expected error for empty customerID, got nil")
	}
}

func TestOrderCreate_ShouldReturnError_WhenItemsIsEmpty(t *testing.T) {
	repo := &mockOrderRepo{}
	producer := &mockEventProducer{}
	svc := NewOrderService(repo, producer)

	_, err := svc.Create(context.Background(), CreateOrderRequest{
		CustomerID: "cust-001",
		Items:      []OrderItemRequest{},
	})

	if err == nil {
		t.Fatal("expected error for empty items, got nil")
	}
}

func TestOrderCreate_ShouldReturnError_WhenItemProductIDIsEmpty(t *testing.T) {
	repo := &mockOrderRepo{}
	producer := &mockEventProducer{}
	svc := NewOrderService(repo, producer)

	_, err := svc.Create(context.Background(), CreateOrderRequest{
		CustomerID: "cust-001",
		Items:      []OrderItemRequest{{ProductID: "", Quantity: 1}},
	})

	if err == nil {
		t.Fatal("expected error for empty productID in item, got nil")
	}
}

func TestOrderCreate_ShouldReturnError_WhenItemQuantityIsZero(t *testing.T) {
	repo := &mockOrderRepo{}
	producer := &mockEventProducer{}
	svc := NewOrderService(repo, producer)

	_, err := svc.Create(context.Background(), CreateOrderRequest{
		CustomerID: "cust-001",
		Items:      []OrderItemRequest{{ProductID: "prod-001", Quantity: 0}},
	})

	if err == nil {
		t.Fatal("expected error for zero quantity, got nil")
	}
}

func TestOrderCreate_ShouldReturnError_WhenRepoSaveFails(t *testing.T) {
	repo := &mockOrderRepo{
		saveFn: func(_ context.Context, _ *order.Order) error {
			return errors.New("db error")
		},
	}
	producer := &mockEventProducer{}
	svc := NewOrderService(repo, producer)

	_, err := svc.Create(context.Background(), validCreateRequest())

	if err == nil {
		t.Fatal("expected error when repo save fails, got nil")
	}
}

func TestOrderCreate_ShouldReturnError_WhenEventPublishFails(t *testing.T) {
	repo := &mockOrderRepo{
		saveFn: func(_ context.Context, _ *order.Order) error { return nil },
	}
	producer := &mockEventProducer{
		publishFn: func(_ context.Context, _ *order.Order) error {
			return errors.New("kafka error")
		},
	}
	svc := NewOrderService(repo, producer)

	_, err := svc.Create(context.Background(), validCreateRequest())

	if err == nil {
		t.Fatal("expected error when event publish fails, got nil")
	}
}

// --- GetByID ---

func TestOrderGetByID_ShouldReturnOrder_WhenFound(t *testing.T) {
	o := newValidOrder(t)
	repo := &mockOrderRepo{
		findByIDFn: func(_ context.Context, _ string) (*order.Order, error) { return o, nil },
	}
	svc := NewOrderService(repo, &mockEventProducer{})

	resp, err := svc.GetByID(context.Background(), "order-001")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.ID != "order-001" {
		t.Errorf("expected ID order-001, got %s", resp.ID)
	}
}

func TestOrderGetByID_ShouldReturnError_WhenNotFound(t *testing.T) {
	repo := &mockOrderRepo{
		findByIDFn: func(_ context.Context, _ string) (*order.Order, error) { return nil, nil },
	}
	svc := NewOrderService(repo, &mockEventProducer{})

	_, err := svc.GetByID(context.Background(), "order-999")

	if err == nil {
		t.Fatal("expected error when order not found, got nil")
	}
}

func TestOrderGetByID_ShouldReturnError_WhenRepoFails(t *testing.T) {
	repo := &mockOrderRepo{
		findByIDFn: func(_ context.Context, _ string) (*order.Order, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewOrderService(repo, &mockEventProducer{})

	_, err := svc.GetByID(context.Background(), "order-001")

	if err == nil {
		t.Fatal("expected error when repo fails, got nil")
	}
}

// --- ListAll ---

func TestOrderListAll_ShouldReturnAllOrders(t *testing.T) {
	o1 := newValidOrder(t)
	o2, _ := order.NewOrder("order-002", "cust-002", []order.OrderItem{
		{ProductID: "prod-002", Quantity: 1},
	})
	repo := &mockOrderRepo{
		findAllFn: func(_ context.Context) ([]*order.Order, error) {
			return []*order.Order{o1, o2}, nil
		},
	}
	svc := NewOrderService(repo, &mockEventProducer{})

	resp, err := svc.ListAll(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(resp) != 2 {
		t.Errorf("expected 2 orders, got %d", len(resp))
	}
}

func TestOrderListAll_ShouldReturnEmptySlice_WhenNoOrders(t *testing.T) {
	repo := &mockOrderRepo{
		findAllFn: func(_ context.Context) ([]*order.Order, error) {
			return []*order.Order{}, nil
		},
	}
	svc := NewOrderService(repo, &mockEventProducer{})

	resp, err := svc.ListAll(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(resp) != 0 {
		t.Errorf("expected empty slice, got %d items", len(resp))
	}
}

func TestOrderListAll_ShouldReturnError_WhenRepoFails(t *testing.T) {
	repo := &mockOrderRepo{
		findAllFn: func(_ context.Context) ([]*order.Order, error) {
			return nil, errors.New("db error")
		},
	}
	svc := NewOrderService(repo, &mockEventProducer{})

	_, err := svc.ListAll(context.Background())

	if err == nil {
		t.Fatal("expected error when repo fails, got nil")
	}
}
