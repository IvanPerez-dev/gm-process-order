package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	application "github.com/ivanperez-dev/gm-order-api/internal/application/product"
)

type ProductHandler struct {
	createUseCase application.CreateProductUseCase
	getUseCase    application.GetProductUseCase
	listUseCase   application.ListProductsUseCase
	updateUseCase application.UpdateProductUseCase
	deleteUseCase application.DeleteProductUseCase
}

func NewProductHandler(
	createUseCase application.CreateProductUseCase,
	getUseCase application.GetProductUseCase,
	listUseCase application.ListProductsUseCase,
	updateUseCase application.UpdateProductUseCase,
	deleteUseCase application.DeleteProductUseCase,
) *ProductHandler {
	return &ProductHandler{
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
		listUseCase:   listUseCase,
		updateUseCase: updateUseCase,
		deleteUseCase: deleteUseCase,
	}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req application.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		NewProblemDetails(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	response, err := h.createUseCase.Create(c.Request.Context(), req)
	if err != nil {
		NewProblemDetails(c, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	response, err := h.getUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		NewProblemDetails(c, http.StatusNotFound, "not found", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) ListAll(c *gin.Context) {
	ids := c.QueryArray("ids")

	products, err := h.listUseCase.ListAll(c.Request.Context(), ids)
	if err != nil {
		NewProblemDetails(c, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req application.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		NewProblemDetails(c, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	response, err := h.updateUseCase.Update(c.Request.Context(), id, req)
	if err != nil {
		NewProblemDetails(c, http.StatusNotFound, "not found", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.deleteUseCase.Delete(c.Request.Context(), id); err != nil {
		NewProblemDetails(c, http.StatusNotFound, "not found", err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
