package order

import "errors"

var (
	// ErrOrderNotFound 订单未找到
	ErrOrderNotFound = errors.New("order not found")

	// ErrInvalidOrderStatus 无效的订单状态
	ErrInvalidOrderStatus = errors.New("invalid order status")

	// ErrOrderAlreadyPaid 订单已支付
	ErrOrderAlreadyPaid = errors.New("order already paid")

	// ErrOrderNotPaid 订单未支付
	ErrOrderNotPaid = errors.New("order not paid")

	// ErrCannotCancelOrder 无法取消订单
	ErrCannotCancelOrder = errors.New("cannot cancel order")

	// ErrEmptyOrder 空订单
	ErrEmptyOrder = errors.New("order must have at least one item")

	// ErrInvalidQuantity 无效的数量
	ErrInvalidQuantity = errors.New("quantity must be greater than zero")

	// ErrPaymentNotFound 支付未找到
	ErrPaymentNotFound = errors.New("payment not found")

	// ErrPaymentFailed 支付失败
	ErrPaymentFailed = errors.New("payment failed")

	// ErrCannotRefund 无法退款
	ErrCannotRefund = errors.New("cannot refund payment")

	// ErrShipmentNotFound 发货未找到
	ErrShipmentNotFound = errors.New("shipment not found")

	// ErrInvalidShipmentStatus 无效的发货状态
	ErrInvalidShipmentStatus = errors.New("invalid shipment status")

	// ErrInvalidAddress 无效的地址
	ErrInvalidAddress = errors.New("invalid address")

	// ErrDifferentCurrency 不同货币
	ErrDifferentCurrency = errors.New("cannot operate on different currencies")
)
