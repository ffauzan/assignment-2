package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"fga-asg-2/pkg/order"
)

type CreateOrderRequest struct {
	OrderedAt    string
	CustomerName string
	Items        []order.Item
}

type UpdateOrderRequest struct {
	OrderedAt    string
	CustomerName string
	Items        []order.Item
}

type OrderService interface {
	CreateOrder(order.Order) (uint, error)
	GetOrder(uint) (*order.Order, error)
	UpdateOrder(order.Order) error
	GetOrders() ([]order.Order, error)
	DeleteOrder(uint) error
}

type orderHandler struct {
	service OrderService
}

func RegisterOrderHandler(r *gin.Engine, service OrderService) {
	handler := &orderHandler{
		service: service,
	}

	order := r.Group("/orders")
	{
		order.POST("/", handler.CreateOrder)
		order.GET("/", handler.GetOrders)
		order.PUT("/:id", handler.UpdateOrder)
		order.GET("/:id", handler.GetOrder)
		order.DELETE("/:id", handler.DeleteOrder)
	}
}

func (h *orderHandler) CreateOrder(c *gin.Context) {
	// Get request body
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	// Create order
	orderId, err := h.service.CreateOrder(order.Order{
		OrderedAt:    req.OrderedAt,
		CustomerName: req.CustomerName,
		Items:        req.Items,
	})
	if err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Return response
	c.JSON(http.StatusOK, BaseResponse{
		Status:  "success",
		Message: "Order created",
		Data: map[string]interface{}{
			"orderId": orderId,
		},
	})
}

func (h *orderHandler) GetOrder(c *gin.Context) {
	// Get order id from url
	orderId := c.Param("id")
	uintOrderId, err := strconv.ParseUint(orderId, 10, 32)
	if err != nil {
		SendErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	// Get order
	order, err := h.service.GetOrder(uint(uintOrderId))
	if err != nil {
		if err.Error() == "record not found" {
			SendErrorResponse(c, err, http.StatusNotFound)
			return
		}
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Return response
	c.JSON(http.StatusOK, BaseResponse{
		Status:  "success",
		Message: "Order found",
		Data:    order,
	})
}

func (h *orderHandler) UpdateOrder(c *gin.Context) {
	// Get order id from url
	orderId := c.Param("id")
	uintOrderId, err := strconv.ParseUint(orderId, 10, 32)
	if err != nil {
		SendErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	// Get request body
	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	// Update order
	err = h.service.UpdateOrder(order.Order{
		OrderID:      uint(uintOrderId),
		OrderedAt:    req.OrderedAt,
		CustomerName: req.CustomerName,
		Items:        req.Items,
	})

	if err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Return response
	c.JSON(http.StatusOK, BaseResponse{
		Status:  "success",
		Message: "Order updated",
		Data:    nil,
	})
}

func (h *orderHandler) GetOrders(c *gin.Context) {
	// Get orders
	orders, err := h.service.GetOrders()
	if err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Return response
	c.JSON(http.StatusOK, BaseResponse{
		Status:  "success",
		Message: "Successfully retrieved orders",
		Data:    orders,
	})
}

func (h *orderHandler) DeleteOrder(c *gin.Context) {
	// Get order id from url
	orderId := c.Param("id")
	uintOrderId, err := strconv.ParseUint(orderId, 10, 32)
	if err != nil {
		SendErrorResponse(c, err, http.StatusBadRequest)
		return
	}

	// Delete order
	err = h.service.DeleteOrder(uint(uintOrderId))
	if err != nil {
		SendErrorResponse(c, err, http.StatusInternalServerError)
		return
	}

	// Return response
	c.JSON(http.StatusOK, BaseResponse{
		Status:  "success",
		Message: "Order deleted",
		Data:    nil,
	})
}
