package order

import "time"

// OrderDTO 订单DTO
type OrderDTO struct {
	ID          string         `json:"id"`
	UserID      string         `json:"user_id"`
	OrderNumber string         `json:"order_number"`
	Status      string         `json:"status"`
	Items       []*OrderItemDTO `json:"items"`
	TotalAmount MoneyDTO       `json:"total_amount"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// OrderItemDTO 订单项DTO
type OrderItemDTO struct {
	ID          string    `json:"id"`
	ProductID   string    `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	UnitPrice   MoneyDTO  `json:"unit_price"`
	Subtotal    MoneyDTO  `json:"subtotal"`
}

// MoneyDTO 金额DTO
type MoneyDTO struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

// CreateOrderItemRequest 创建订单项请求
type CreateOrderItemRequest struct {
	ProductID   string  `json:"product_id" validate:"required"`
	ProductName string  `json:"product_name" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required,gte=1"`
	UnitPrice   float64 `json:"unit_price" validate:"required,gt=0"`
	Currency    string  `json:"currency" validate:"required,len=3"`
}

// PaymentDTO 支付DTO
type PaymentDTO struct {
	ID              string    `json:"id"`
	OrderID         string    `json:"order_id"`
	Amount          MoneyDTO  `json:"amount"`
	Method          string    `json:"method"`
	Status          string    `json:"status"`
	TransactionID   string    `json:"transaction_id,omitempty"`
	GatewayResponse string    `json:"gateway_response,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ProcessPaymentRequest 处理支付请求
type ProcessPaymentRequest struct {
	Method string `json:"method" validate:"required"`
}

// ShipmentDTO 发货DTO
type ShipmentDTO struct {
	ID             string       `json:"id"`
	OrderID        string       `json:"order_id"`
	TrackingNumber string       `json:"tracking_number"`
	Carrier        string       `json:"carrier"`
	ShippingMethod string       `json:"shipping_method"`
	Address        AddressDTO   `json:"address"`
	Status         string       `json:"status"`
	EstimatedDate  *time.Time   `json:"estimated_date,omitempty"`
	ShippedAt      *time.Time   `json:"shipped_at,omitempty"`
	DeliveredAt    *time.Time   `json:"delivered_at,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// AddressDTO 地址DTO
type AddressDTO struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// CreateShipmentRequest 创建发货请求
type CreateShipmentRequest struct {
	ShippingMethod string     `json:"shipping_method" validate:"required"`
	Address        AddressDTO `json:"address" validate:"required"`
}

// UpdateShipmentRequest 更新发货请求
type UpdateShipmentRequest struct {
	TrackingNumber string `json:"tracking_number" validate:"required"`
	Carrier        string `json:"carrier" validate:"required"`
}

// ListOrdersRequest 列出订单请求
type ListOrdersRequest struct {
	Page     int `json:"page" validate:"gte=1"`
	PageSize int `json:"page_size" validate:"gte=1,lte=100"`
}

// ListOrdersResponse 列出订单响应
type ListOrdersResponse struct {
	Orders     []*OrderDTO `json:"orders"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

