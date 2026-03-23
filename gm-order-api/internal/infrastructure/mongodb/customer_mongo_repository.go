package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/ivanperez-dev/gm-order-api/internal/domain/customer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type customerDoc struct {
	ID        string    `bson:"_id"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	IsActieve bool      `bson:"isActive"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type CustomerMongoRepository struct {
	collection *mongo.Collection
}

func NewCustomerMongoRepository(db *mongo.Database) *CustomerMongoRepository {
	collection := db.Collection("customers")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetName("idx_email").SetUnique(true),
	})

	return &CustomerMongoRepository{collection: collection}
}

func (r *CustomerMongoRepository) Save(ctx context.Context, c *customer.Customer) error {
	doc := customerToDoc(c)
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return errors.New("error saving customer: " + err.Error())
	}
	return nil
}

func (r *CustomerMongoRepository) FindByID(ctx context.Context, id string) (*customer.Customer, error) {
	var doc customerDoc
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, errors.New("error finding customer: " + err.Error())
	}
	return customerToDomain(&doc), nil
}

func (r *CustomerMongoRepository) FindAll(ctx context.Context) ([]*customer.Customer, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.New("error listing customers: " + err.Error())
	}
	defer cursor.Close(ctx)

	var customers []*customer.Customer
	for cursor.Next(ctx) {
		var doc customerDoc
		if err := cursor.Decode(&doc); err != nil {
			return nil, errors.New("error decoding customer: " + err.Error())
		}
		customers = append(customers, customerToDomain(&doc))
	}
	return customers, nil
}

func (r *CustomerMongoRepository) Update(ctx context.Context, c *customer.Customer) error {
	doc := customerToDoc(c)
	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": c.ID}, doc)
	if err != nil {
		return errors.New("error updating customer: " + err.Error())
	}
	if result.MatchedCount == 0 {
		return errors.New("customer not found: " + c.ID)
	}
	return nil
}

func (r *CustomerMongoRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.New("error deleting customer: " + err.Error())
	}
	if result.DeletedCount == 0 {
		return errors.New("customer not found: " + id)
	}
	return nil
}

func customerToDoc(c *customer.Customer) customerDoc {
	return customerDoc{
		ID:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		IsActieve: c.IsActive,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func customerToDomain(doc *customerDoc) *customer.Customer {
	return &customer.Customer{
		ID:        doc.ID,
		Name:      doc.Name,
		Email:     doc.Email,
		IsActive:  doc.IsActieve,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}
