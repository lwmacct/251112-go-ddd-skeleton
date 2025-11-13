package order

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Service 订单领域服务
type Service struct {
	orderRepo    OrderRepository
	paymentRepo  PaymentRepository
	shipmentRepo ShipmentRepository
}

// NewService 创建订单领域服务
func NewService(orderRepo OrderRepository, paymentRepo PaymentRepository, shipmentRepo ShipmentRepository) *Service {
	return &Service{
		orderRepo:    orderRepo,
		paymentRepo:  paymentRepo,
		shipmentRepo: shipmentRepo,
	}
}

// GenerateOrderNumber 生成订单号
func (s *Service) GenerateOrderNumber() string {
	return fmt.Sprintf("ORD-%d", time.Now().UnixNano())
}

// ValidateOrderForPayment 验证订单是否可以支付
func (s *Service) ValidateOrderForPayment(ctx context.Context, orderID string) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != StatusPending {
		return errors.New("only pending orders can be paid")
	}

	if len(order.Items) == 0 {
		return errors.New("order must have at least one item")
	}

	if !order.TotalAmount.IsPositive() {
		return errors.New("order total must be positive")
	}

	return nil
}

// ValidateOrderForShipment 验证订单是否可以发货
func (s *Service) ValidateOrderForShipment(ctx context.Context, orderID string) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != StatusPaid {
		return errors.New("only paid orders can be shipped")
	}

	return nil
}

// CanCancelOrder 验证订单是否可以取消
func (s *Service) CanCancelOrder(ctx context.Context, orderID string) (bool, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return false, err
	}

	return order.CanBeCancelled(), nil
}

// ValidateRefund 验证是否可以退款
func (s *Service) ValidateRefund(ctx context.Context, orderID string) error {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != StatusPaid && order.Status != StatusCompleted {
		return errors.New("only paid or completed orders can be refunded")
	}

	payment, err := s.paymentRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	if !payment.CanBeRefunded() {
		return errors.New("payment cannot be refunded")
	}

	return nil
}

// CalculateEstimatedDelivery 计算预计送达时间（示例：7天后）
func (s *Service) CalculateEstimatedDelivery(shippingMethod string) time.Time {
	daysToAdd := 7 // 默认7天

	switch shippingMethod {
	case "express":
		daysToAdd = 2
	case "standard":
		daysToAdd = 5
	case "economy":
		daysToAdd = 10
	}

	return time.Now().AddDate(0, 0, daysToAdd)
}

