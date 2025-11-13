package order

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/adapters/http/response"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/application/order"
)

// Handler 订单处理器
type Handler struct {
	orderService *order.Service
}

// NewHandler 创建订单处理器
func NewHandler(orderService *order.Service) *Handler {
	return &Handler{
		orderService: orderService,
	}
}

// CreateOrder 创建订单
func (h *Handler) CreateOrder(c *gin.Context) {
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
func (h *Handler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")

	dto, err := h.orderService.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto)
}

// ListOrders 列出订单
func (h *Handler) ListOrders(c *gin.Context) {
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
func (h *Handler) CancelOrder(c *gin.Context) {
	orderID := c.Param("id")

	if err := h.orderService.CancelOrder(c.Request.Context(), orderID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Order cancelled successfully"})
}

// ProcessPayment 处理支付
func (h *Handler) ProcessPayment(c *gin.Context) {
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
func (h *Handler) CreateShipment(c *gin.Context) {
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

// GetPayment 获取支付信息
func (h *Handler) GetPayment(c *gin.Context) {
	orderID := c.Param("id")

	// TODO: 实现获取支付信息逻辑
	_ = orderID
	response.Success(c, gin.H{"order_id": orderID, "status": "pending"})
}

// GetShipment 获取发货信息
func (h *Handler) GetShipment(c *gin.Context) {
	orderID := c.Param("id")

	// TODO: 实现获取发货信息逻辑
	_ = orderID
	response.Success(c, gin.H{"order_id": orderID, "status": "preparing"})
}

// ========== 管理员订单管理端点（需要admin权限）==========

// ListAllOrders 列出所有订单
// GET /api/admin/orders
func (h *Handler) ListAllOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	req := order.ListOrdersRequest{
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := h.orderService.ListOrders(c.Request.Context(), "", req) // 空 userID 表示所有订单
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, resp)
}

// GetOrderByID 获取指定订单详情（管理员用）
// GET /api/admin/orders/:id
func (h *Handler) GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")

	dto, err := h.orderService.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto)
}

// UpdateOrderStatus 更新订单状态
// PUT /api/admin/orders/:id/status
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
		Note   string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// TODO: 实现订单状态更新逻辑
	_ = orderID
	response.Success(c, gin.H{"message": "订单状态已更新", "status": req.Status})
}
