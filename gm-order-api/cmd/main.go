package main

import (
	"context"
	"log"
	"strings"

	customerapplication "github.com/ivanperez-dev/gm-order-api/internal/application/customer"
	orderapplication "github.com/ivanperez-dev/gm-order-api/internal/application/order"
	productapplication "github.com/ivanperez-dev/gm-order-api/internal/application/product"
	infrahttp "github.com/ivanperez-dev/gm-order-api/internal/infrastructure/http"
	"github.com/ivanperez-dev/gm-order-api/internal/infrastructure/kafka"
	"github.com/ivanperez-dev/gm-order-api/internal/infrastructure/mongodb"
)

func main() {

	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoClient, err := connectMongo(mongoURI)
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	db := mongoClient.Database(getEnv("MONGO_DB", "gm_orders"))
	orderRepo := mongodb.NewOrderMongoRepository(db)

	brokers := strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")
	topic := getEnv("KAFKA_TOPIC", "orders.created")
	kafkaWriter := kafka.NewKafkaWriter(brokers, topic)
	producer := kafka.NewKafkaEventProducer(kafkaWriter)
	defer producer.Close()

	orderService := orderapplication.NewOrderService(orderRepo, producer)
	productRepo := mongodb.NewProductMongoRepository(db)
	productService := productapplication.NewProductService(productRepo)
	customerRepo := mongodb.NewCustomerMongoRepository(db)
	customerService := customerapplication.NewCustomerService(customerRepo)

	runSeed(context.Background(), customerRepo, productRepo)

	orderHandler := infrahttp.NewOrderHandler(
		orderService,
		orderService,
		orderService,
	)
	productHandler := infrahttp.NewProductHandler(
		productService,
		productService,
		productService,
		productService,
		productService,
	)
	customerHandler := infrahttp.NewCustomerHandler(
		customerService,
		customerService,
		customerService,
		customerService,
		customerService,
	)
	router := infrahttp.NewRouter(orderHandler, productHandler, customerHandler)

	port := getEnv("PORT", "8081")
	log.Printf("gm-order-api running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
