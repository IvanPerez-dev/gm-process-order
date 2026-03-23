package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/ivanperez-dev/gm-order-api/internal/domain/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderItemDoc struct {
	ProductID string  `bson:"productId"`
	Name      string  `bson:"name"`
	Price     float64 `bson:"price"`
	Quantity  int     `bson:"quantity"`
}

type OrderCustomerDoc struct {
	ID    string `bson:"_id"`
	Name  string `bson:"name"`
	Email string `bson:"email"`
}

type orderDoc struct {
	ID          string           `bson:"_id"`
	Customer    OrderCustomerDoc `bson:"customer"`
	Items       []orderItemDoc   `bson:"products"`
	Status      string           `bson:"status"`
	RetryCount  int              `bson:"retryCount"`
	ErrorMsg    string           `bson:"errorMsg,omitempty"`
	CreatedAt   time.Time        `bson:"createdAt"`
	ProcessedAt *time.Time       `bson:"processedAt,omitempty"`
}

type OrderMongoRepository struct {
	collection *mongo.Collection
}

func NewOrderMongoRepository(db *mongo.Database) *OrderMongoRepository {
	collection := db.Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "status", Value: 1}},
		Options: options.Index().SetName("idx_status"),
	})

	return &OrderMongoRepository{collection: collection}
}

func (r *OrderMongoRepository) Save(ctx context.Context, o *order.Order) error {
	doc := toDoc(o)
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return errors.New("error saving order: " + err.Error())
	}
	return nil
}

func (r *OrderMongoRepository) FindByID(ctx context.Context, id string) (*order.Order, error) {
	var doc orderDoc
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, errors.New("error finding order: " + err.Error())
	}
	return toDomain(&doc), nil
}

func (r *OrderMongoRepository) FindAll(ctx context.Context) ([]*order.Order, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.New("error listing orders: " + err.Error())
	}
	defer cursor.Close(ctx)

	var orders []*order.Order
	for cursor.Next(ctx) {
		var doc orderDoc
		if err := cursor.Decode(&doc); err != nil {
			return nil, errors.New("error decoding order: " + err.Error())
		}
		orders = append(orders, toDomain(&doc))
	}
	return orders, nil
}

func (r *OrderMongoRepository) Update(ctx context.Context, o *order.Order) error {
	doc := toDoc(o)
	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": o.ID}, doc)
	if err != nil {
		return errors.New("error updating order: " + err.Error())
	}
	if result.MatchedCount == 0 {
		return errors.New("order not found: " + o.ID)
	}
	return nil
}

func toDoc(o *order.Order) orderDoc {
	items := make([]orderItemDoc, len(o.Items))
	for i, item := range o.Items {
		items[i] = orderItemDoc{
			ProductID: item.ProductID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
		}
	}
	return orderDoc{
		ID: o.ID,
		Customer: OrderCustomerDoc{
			ID: o.CustomerID,
		},
		Items:       items,
		Status:      string(o.Status),
		RetryCount:  o.RetryCount,
		ErrorMsg:    o.ErrorMsg,
		CreatedAt:   o.CreatedAt,
		ProcessedAt: o.ProcessedAt,
	}
}

func toDomain(doc *orderDoc) *order.Order {
	items := make([]order.OrderItem, len(doc.Items))
	for i, item := range doc.Items {
		items[i] = order.OrderItem{
			ProductID: item.ProductID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
		}
	}
	return &order.Order{
		ID:          doc.ID,
		CustomerID:  doc.Customer.ID,
		Items:       items,
		Status:      order.OrderStatus(doc.Status),
		RetryCount:  doc.RetryCount,
		ErrorMsg:    doc.ErrorMsg,
		CreatedAt:   doc.CreatedAt,
		ProcessedAt: doc.ProcessedAt,
	}
}
