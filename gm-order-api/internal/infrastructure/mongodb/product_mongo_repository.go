package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/ivanperez-dev/gm-order-api/internal/domain/product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productDoc struct {
	ID        string    `bson:"_id"`
	Name      string    `bson:"name"`
	Price     float64   `bson:"price"`
	Stock     int       `bson:"stock"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type ProductMongoRepository struct {
	collection *mongo.Collection
}

func NewProductMongoRepository(db *mongo.Database) *ProductMongoRepository {
	collection := db.Collection("products")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetName("idx_name"),
	})

	return &ProductMongoRepository{collection: collection}
}

func (r *ProductMongoRepository) Save(ctx context.Context, p *product.Product) error {
	doc := productToDoc(p)
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return errors.New("error saving product: " + err.Error())
	}
	return nil
}

func (r *ProductMongoRepository) FindByID(ctx context.Context, id string) (*product.Product, error) {
	var doc productDoc
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, errors.New("error finding product: " + err.Error())
	}
	return productToDomain(&doc), nil
}

func (r *ProductMongoRepository) FindAll(ctx context.Context) ([]*product.Product, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.New("error listing products: " + err.Error())
	}
	defer cursor.Close(ctx)

	var products []*product.Product
	for cursor.Next(ctx) {
		var doc productDoc
		if err := cursor.Decode(&doc); err != nil {
			return nil, errors.New("error decoding product: " + err.Error())
		}
		products = append(products, productToDomain(&doc))
	}
	return products, nil
}

func (r *ProductMongoRepository) FindByIDs(ctx context.Context, ids []string) ([]*product.Product, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, errors.New("error listing products by ids: " + err.Error())
	}
	defer cursor.Close(ctx)

	var products []*product.Product
	for cursor.Next(ctx) {
		var doc productDoc
		if err := cursor.Decode(&doc); err != nil {
			return nil, errors.New("error decoding product: " + err.Error())
		}
		products = append(products, productToDomain(&doc))
	}
	return products, nil
}

func (r *ProductMongoRepository) Update(ctx context.Context, p *product.Product) error {
	doc := productToDoc(p)
	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": p.ID}, doc)
	if err != nil {
		return errors.New("error updating product: " + err.Error())
	}
	if result.MatchedCount == 0 {
		return errors.New("product not found: " + p.ID)
	}
	return nil
}

func (r *ProductMongoRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.New("error deleting product: " + err.Error())
	}
	if result.DeletedCount == 0 {
		return errors.New("product not found: " + id)
	}
	return nil
}

func productToDoc(p *product.Product) productDoc {
	return productDoc{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func productToDomain(doc *productDoc) *product.Product {
	return &product.Product{
		ID:        doc.ID,
		Name:      doc.Name,
		Price:     doc.Price,
		Stock:     doc.Stock,
		CreatedAt: doc.CreatedAt,
		UpdatedAt: doc.UpdatedAt,
	}
}
