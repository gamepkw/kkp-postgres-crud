package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	constant "postgres-crud/app/internal/constants"
	model "postgres-crud/app/internal/models"
	service "postgres-crud/app/internal/services"

	"github.com/labstack/echo/v4"
)

type orderHandler struct {
	service service.IOrderService
}

type IOrderHandler interface {
	CreateOrder(c echo.Context) error
	GetAll(c echo.Context) error
	GetByID(c echo.Context) error
	UpdateOrderStatus(c echo.Context) error
}

func NewOrderHandler(service service.IOrderService) *orderHandler {
	return &orderHandler{service: service}
}

func (h *orderHandler) CreateOrder(c echo.Context) error {
	var order model.Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	if err := c.Validate(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	if err := h.service.Create(ctx, &order); err != nil {
		if err == context.DeadlineExceeded {
			return c.JSON(http.StatusRequestTimeout, map[string]string{"error": "Request timed out"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create order: " + err.Error()})
	}
	return c.JSON(http.StatusCreated, constant.RespSuccessWithData(order))
}

func (h *orderHandler) GetAll(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()

	page, _ := strconv.Atoi(c.QueryParam("page"))
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))

	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 100
	}

	orders, err := h.service.GetAll(ctx, page, perPage)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders"})
	}
	return c.JSON(http.StatusOK, constant.RespSuccessWithData(orders))
}

func (h *orderHandler) GetByID(c echo.Context) error {
	idParam := c.Param("order_id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	order, err := h.service.GetByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders"})
	}
	return c.JSON(http.StatusOK, constant.RespSuccessWithData(order))
}

func (h *orderHandler) UpdateOrderStatus(c echo.Context) error {
	idParam := c.Param("order_id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	var order model.Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	if err := c.Validate(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*10)
	defer cancel()
	if err := h.service.Update(ctx, id, order.Status); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update order"})
	}
	return c.JSON(http.StatusOK, constant.RespSuccess())
}
