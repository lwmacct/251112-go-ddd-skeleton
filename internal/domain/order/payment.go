package order

import (
	"errors"
	"time"
)

// PaymentStatus 支付状态
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// PaymentMethod 支付方式值对象
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodDebitCard  PaymentMethod = "debit_card"
	PaymentMethodPayPal     PaymentMethod = "paypal"
	PaymentMethodStripe     PaymentMethod = "stripe"
	PaymentMethodCash       PaymentMethod = "cash"
)

// Payment 支付实体
type Payment struct {
	ID              string
	OrderID         string
	Amount          Money
	Method          PaymentMethod
	Status          PaymentStatus
	TransactionID   string
	GatewayResponse string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewPayment 创建支付
func NewPayment(orderID string, amount Money, method PaymentMethod) (*Payment, error) {
	if orderID == "" {
		return nil, errors.New("orderID cannot be empty")
	}
	if !amount.IsPositive() {
		return nil, errors.New("amount must be positive")
	}

	return &Payment{
		OrderID:   orderID,
		Amount:    amount,
		Method:    method,
		Status:    PaymentStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// MarkAsCompleted 标记为已完成
func (p *Payment) MarkAsCompleted(transactionID, gatewayResponse string) error {
	if p.Status != PaymentStatusPending {
		return errors.New("can only complete pending payments")
	}
	p.Status = PaymentStatusCompleted
	p.TransactionID = transactionID
	p.GatewayResponse = gatewayResponse
	p.UpdatedAt = time.Now()
	return nil
}

// MarkAsFailed 标记为失败
func (p *Payment) MarkAsFailed(gatewayResponse string) error {
	if p.Status == PaymentStatusCompleted || p.Status == PaymentStatusRefunded {
		return errors.New("cannot mark completed or refunded payment as failed")
	}
	p.Status = PaymentStatusFailed
	p.GatewayResponse = gatewayResponse
	p.UpdatedAt = time.Now()
	return nil
}

// Refund 退款
func (p *Payment) Refund() error {
	if p.Status != PaymentStatusCompleted {
		return errors.New("can only refund completed payments")
	}
	p.Status = PaymentStatusRefunded
	p.UpdatedAt = time.Now()
	return nil
}

// IsCompleted 是否已完成
func (p *Payment) IsCompleted() bool {
	return p.Status == PaymentStatusCompleted
}

// CanBeRefunded 是否可以退款
func (p *Payment) CanBeRefunded() bool {
	return p.Status == PaymentStatusCompleted
}
