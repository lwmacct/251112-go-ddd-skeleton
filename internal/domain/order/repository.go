package order

import "context"

// OrderRepository 订单仓储接口
type OrderRepository interface {
	// Create 创建订单
	Create(ctx context.Context, order *Order) error

	// Update 更新订单
	Update(ctx context.Context, order *Order) error

	// FindByID 根据ID查找订单
	FindByID(ctx context.Context, id string) (*Order, error)

	// FindByOrderNumber 根据订单号查找订单
	FindByOrderNumber(ctx context.Context, orderNumber string) (*Order, error)

	// ListByUserID 列出用户的订单
	ListByUserID(ctx context.Context, userID string, offset, limit int) ([]*Order, int64, error)

	// Delete 删除订单
	Delete(ctx context.Context, id string) error
}

// PaymentRepository 支付仓储接口
type PaymentRepository interface {
	// Create 创建支付
	Create(ctx context.Context, payment *Payment) error

	// Update 更新支付
	Update(ctx context.Context, payment *Payment) error

	// FindByID 根据ID查找支付
	FindByID(ctx context.Context, id string) (*Payment, error)

	// FindByOrderID 根据订单ID查找支付
	FindByOrderID(ctx context.Context, orderID string) (*Payment, error)

	// FindByTransactionID 根据交易ID查找支付
	FindByTransactionID(ctx context.Context, transactionID string) (*Payment, error)
}

// ShipmentRepository 发货仓储接口
type ShipmentRepository interface {
	// Create 创建发货
	Create(ctx context.Context, shipment *Shipment) error

	// Update 更新发货
	Update(ctx context.Context, shipment *Shipment) error

	// FindByID 根据ID查找发货
	FindByID(ctx context.Context, id string) (*Shipment, error)

	// FindByOrderID 根据订单ID查找发货
	FindByOrderID(ctx context.Context, orderID string) (*Shipment, error)

	// FindByTrackingNumber 根据追踪号查找发货
	FindByTrackingNumber(ctx context.Context, trackingNumber string) (*Shipment, error)
}
