package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-ddd-skeleton/internal/adapters/http/response"
	"github.com/yourusername/go-ddd-skeleton/internal/application/order"
)

// OrderHandler 订单处理器
type OrderHandler struct {
	orderService *order.Service
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(orderService *order.Service) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder 创建订单
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetString("userID")

	var req order.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	dto, err := h.orderService.CreateOrder(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, dto)
}

// GetOrder 获取订单
func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")

	dto, err := h.orderService.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto)
}

// ListOrders 列出订单
func (h *OrderHandler) ListOrders(c *gin.Context) {
	userID := c.GetString("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	req := order.ListOrdersRequest{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := h.orderService.ListOrders(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, resp)
}

// CancelOrder 取消订单
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")

	if err := h.orderService.CancelOrder(c.Request.Context(), orderID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Order cancelled successfully"})
}

// ProcessPayment 处理支付
func (h *OrderHandler) ProcessPayment(c *gin.Context) {
	orderID := c.Param("id")

	var req order.ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	dto, err := h.orderService.ProcessPayment(c.Request.Context(), orderID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto)
}

// CreateShipment 创建发货
func (h *OrderHandler) CreateShipment(c *gin.Context) {
	orderID := c.Param("id")

	var req order.CreateShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	dto, err := h.orderService.CreateShipment(c.Request.Context(), orderID, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, dto)
}
