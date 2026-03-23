package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	application "github.com/ivanperez-dev/gm-order-api/internal/application/order"
)

type OrderHandler struct {
	createUseCase application.CreateOrderUseCase
	getUseCase    application.GetOrderUseCase
	listUseCase   application.ListOrdersUseCase
}

func NewOrderHandler(
	createUseCase application.CreateOrderUseCase,
	getUseCase application.GetOrderUseCase,
	listUseCase application.ListOrdersUseCase,
) *OrderHandler {
	return &OrderHandler{
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
		listUseCase:   listUseCase,
	}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req application.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		NewProblemDetails(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	response, err := h.createUseCase.Create(c.Request.Context(), req)
	if err != nil {

		NewProblemDetails(c, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	c.Header("Location", "/api/v1/orders/"+response.ID)
	c.JSON(http.StatusCreated, response)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	response, err := h.getUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		NewProblemDetails(c, http.StatusNotFound, "not found", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) ListAll(c *gin.Context) {
	orders, err := h.listUseCase.ListAll(c.Request.Context())
	if err != nil {

		NewProblemDetails(c, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	c.JSON(http.StatusOK, orders)
}
