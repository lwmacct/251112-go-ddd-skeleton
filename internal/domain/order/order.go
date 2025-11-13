package order

import (
	"errors"
	"time"
)

// OrderStatus 订单状态值对象
type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPaid      OrderStatus = "paid"
	StatusCancelled OrderStatus = "cancelled"
	StatusCompleted OrderStatus = "completed"
	StatusRefunded  OrderStatus = "refunded"
)

// IsValid 检查订单状态是否有效
func (s OrderStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusPaid, StatusCancelled, StatusCompleted, StatusRefunded:
		return true
	}
	return false
}

// Order 订单聚合根
type Order struct {
	ID          string
	UserID      string
	OrderNumber string
	Status      OrderStatus
	Items       []*OrderItem
	TotalAmount Money
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewOrder 创建新订单
func NewOrder(userID, orderNumber string) (*Order, error) {
	if userID == "" {
		return nil, errors.New("userID cannot be empty")
	}
	if orderNumber == "" {
		return nil, errors.New("orderNumber cannot be empty")
	}

	return &Order{
		UserID:      userID,
		OrderNumber: orderNumber,
		Status:      StatusPending,
		Items:       make([]*OrderItem, 0),
		TotalAmount: NewMoney(0, "USD"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// AddItem 添加订单项
func (o *Order) AddItem(item *OrderItem) error {
	if item == nil {
		return errors.New("item cannot be nil")
	}

	o.Items = append(o.Items, item)
	o.calculateTotal()
	o.UpdatedAt = time.Now()
	return nil
}

// RemoveItem 移除订单项
func (o *Order) RemoveItem(itemID string) error {
	for i, item := range o.Items {
		if item.ID == itemID {
			o.Items = append(o.Items[:i], o.Items[i+1:]...)
			o.calculateTotal()
			o.UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("item not found")
}

// calculateTotal 计算订单总额
func (o *Order) calculateTotal() {
	total := 0.0
	currency := "USD"

	for _, item := range o.Items {
		total += item.Subtotal.Amount
		currency = item.Subtotal.Currency
	}

	o.TotalAmount = NewMoney(total, currency)
}

// MarkAsPaid 标记为已支付
func (o *Order) MarkAsPaid() error {
	if o.Status != StatusPending {
		return errors.New("can only mark pending orders as paid")
	}
	o.Status = StatusPaid
	o.UpdatedAt = time.Now()
	return nil
}

// Cancel 取消订单
func (o *Order) Cancel() error {
	if o.Status == StatusCompleted || o.Status == StatusRefunded {
		return errors.New("cannot cancel completed or refunded orders")
	}
	o.Status = StatusCancelled
	o.UpdatedAt = time.Now()
	return nil
}

// Complete 完成订单
func (o *Order) Complete() error {
	if o.Status != StatusPaid {
		return errors.New("can only complete paid orders")
	}
	o.Status = StatusCompleted
	o.UpdatedAt = time.Now()
	return nil
}

// Refund 退款订单
func (o *Order) Refund() error {
	if o.Status != StatusPaid && o.Status != StatusCompleted {
		return errors.New("can only refund paid or completed orders")
	}
	o.Status = StatusRefunded
	o.UpdatedAt = time.Now()
	return nil
}

// CanBeCancelled 是否可以取消
func (o *Order) CanBeCancelled() bool {
	return o.Status != StatusCompleted && o.Status != StatusRefunded
}

// OrderItem 订单项实体
type OrderItem struct {
	ID          string
	OrderID     string
	ProductID   string
	ProductName string
	Quantity    int
	UnitPrice   Money
	Subtotal    Money
	CreatedAt   time.Time
}

// NewOrderItem 创建订单项
func NewOrderItem(orderID, productID, productName string, quantity int, unitPrice Money) (*OrderItem, error) {
	if orderID == "" {
		return nil, errors.New("orderID cannot be empty")
	}
	if productID == "" {
		return nil, errors.New("productID cannot be empty")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	subtotal := NewMoney(unitPrice.Amount*float64(quantity), unitPrice.Currency)

	return &OrderItem{
		OrderID:     orderID,
		ProductID:   productID,
		ProductName: productName,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		Subtotal:    subtotal,
		CreatedAt:   time.Now(),
	}, nil
}

// UpdateQuantity 更新数量
func (oi *OrderItem) UpdateQuantity(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	oi.Quantity = quantity
	oi.Subtotal = NewMoney(oi.UnitPrice.Amount*float64(quantity), oi.UnitPrice.Currency)
	return nil
}

// Money 金额值对象
type Money struct {
	Amount   float64
	Currency string
}

// NewMoney 创建金额值对象
func NewMoney(amount float64, currency string) Money {
	return Money{
		Amount:   amount,
		Currency: currency,
	}
}

// Add 加法
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("cannot add different currencies")
	}
	return NewMoney(m.Amount+other.Amount, m.Currency), nil
}

// Subtract 减法
func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("cannot subtract different currencies")
	}
	return NewMoney(m.Amount-other.Amount, m.Currency), nil
}

// Multiply 乘法
func (m Money) Multiply(factor float64) Money {
	return NewMoney(m.Amount*factor, m.Currency)
}

// IsZero 是否为零
func (m Money) IsZero() bool {
	return m.Amount == 0
}

// IsPositive 是否为正数
func (m Money) IsPositive() bool {
	return m.Amount > 0
}
