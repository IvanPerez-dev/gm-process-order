package http

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(orderHandler *OrderHandler, productHandler *ProductHandler, customerHandler *CustomerHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		orders := v1.Group("/orders")
		{
			orders.POST("", orderHandler.Create)
			orders.GET("", orderHandler.ListAll)
			orders.GET("/:id", orderHandler.GetByID)
		}

		products := v1.Group("/products")
		{
			products.POST("", productHandler.Create)
			products.GET("", productHandler.ListAll)
			products.GET("/:id", productHandler.GetByID)
			products.PUT("/:id", productHandler.Update)
			products.DELETE("/:id", productHandler.Delete)
		}

		customers := v1.Group("/customers")
		{
			customers.POST("", customerHandler.Create)
			customers.GET("", customerHandler.ListAll)
			customers.GET("/:id", customerHandler.GetByID)
			customers.PUT("/:id", customerHandler.Update)
			customers.DELETE("/:id", customerHandler.Delete)
		}
	}

	return r
}
