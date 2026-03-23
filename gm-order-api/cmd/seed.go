package main

import (
	"context"
	"log"

	customerdomain "github.com/ivanperez-dev/gm-order-api/internal/domain/customer"
	productdomain "github.com/ivanperez-dev/gm-order-api/internal/domain/product"
)

var seedCustomers = []struct {
	id    string
	name  string
	email string
}{
	{"c1a2b3c4-0001-4000-a000-000000000001", "Juan Pérez", "juan.perez@grupomariposa.com"},
	{"c1a2b3c4-0001-4000-a000-000000000002", "María García", "maria.garcia@grupomariposa.com"},
	{"c1a2b3c4-0001-4000-a000-000000000003", "Carlos López", "carlos.lopez@grupomariposa.com"},
}

var seedProducts = []struct {
	id    string
	name  string
	price float64
	stock int
}{
	{"p1a2b3c4-0001-4000-a000-000000000001", "Producto Alpha", 99.99, 100},
	{"p1a2b3c4-0001-4000-a000-000000000002", "Producto Beta", 149.50, 50},
	{"p1a2b3c4-0001-4000-a000-000000000003", "Producto Gamma", 49.00, 200},
	{"p1a2b3c4-0001-4000-a000-000000000004", "Producto Delta", 299.99, 25},
}

func runSeed(ctx context.Context, customerRepo customerdomain.Repository, productRepo productdomain.Repository) {
	seedCustomerData(ctx, customerRepo)
	seedProductData(ctx, productRepo)
}

func seedCustomerData(ctx context.Context, repo customerdomain.Repository) {
	for _, s := range seedCustomers {
		existing, err := repo.FindByID(ctx, s.id)
		if err != nil {
			log.Printf("[seed] error checking customer %s: %v", s.id, err)
			continue
		}
		if existing != nil {
			log.Printf("[seed] customer already exists: %s", s.id)
			continue
		}

		c, err := customerdomain.NewCustomer(s.id, s.name, s.email)
		if err != nil {
			log.Printf("[seed] invalid customer data for %s: %v", s.id, err)
			continue
		}
		if err := repo.Save(ctx, c); err != nil {
			log.Printf("[seed] error saving customer %s: %v", s.id, err)
			continue
		}
		log.Printf("[seed] customer created: %s (%s)", s.name, s.id)
	}
}

func seedProductData(ctx context.Context, repo productdomain.Repository) {
	for _, s := range seedProducts {
		existing, err := repo.FindByID(ctx, s.id)
		if err != nil {
			log.Printf("[seed] error checking product %s: %v", s.id, err)
			continue
		}
		if existing != nil {
			log.Printf("[seed] product already exists: %s", s.id)
			continue
		}

		p, err := productdomain.NewProduct(s.id, s.name, s.price, s.stock)
		if err != nil {
			log.Printf("[seed] invalid product data for %s: %v", s.id, err)
			continue
		}
		if err := repo.Save(ctx, p); err != nil {
			log.Printf("[seed] error saving product %s: %v", s.id, err)
			continue
		}
		log.Printf("[seed] product created: %s (%s)", s.name, s.id)
	}
}
