package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	application "github.com/ivanperez-dev/gm-order-api/internal/application/customer"
)

type CustomerHandler struct {
	createUseCase application.CreateCustomerUseCase
	getUseCase    application.GetCustomerUseCase
	listUseCase   application.ListCustomersUseCase
	updateUseCase application.UpdateCustomerUseCase
	deleteUseCase application.DeleteCustomerUseCase
}

func NewCustomerHandler(
	createUseCase application.CreateCustomerUseCase,
	getUseCase application.GetCustomerUseCase,
	listUseCase application.ListCustomersUseCase,
	updateUseCase application.UpdateCustomerUseCase,
	deleteUseCase application.DeleteCustomerUseCase,
) *CustomerHandler {
	return &CustomerHandler{
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
		listUseCase:   listUseCase,
		updateUseCase: updateUseCase,
		deleteUseCase: deleteUseCase,
	}
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var req application.CreateCustomerRequest
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

func (h *CustomerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	response, err := h.getUseCase.GetByID(c.Request.Context(), id)
	if err != nil {

		NewProblemDetails(c, http.StatusNotFound, "not found", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *CustomerHandler) ListAll(c *gin.Context) {
	customers, err := h.listUseCase.ListAll(c.Request.Context())
	if err != nil {

		NewProblemDetails(c, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req application.UpdateCustomerRequest
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

func (h *CustomerHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.deleteUseCase.Delete(c.Request.Context(), id); err != nil {
		NewProblemDetails(c, http.StatusNotFound, "not found", err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
