package order

import (
	"context"
	"math/rand"
	"time"

	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/order"
	"github.com/oklog/ulid/v2"
)

// CreateOrder 创建订单（命令）
func (s *Service) CreateOrder(ctx context.Context, userID string, req CreateOrderRequest) (*OrderDTO, error) {
	// 生成订单号
	orderNumber := s.orderService.GenerateOrderNumber()

	// 创建订单
	o, err := order.NewOrder(userID, orderNumber)
	if err != nil {
		return nil, err
	}

	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	o.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	// 添加订单项
	for _, itemReq := range req.Items {
		itemID := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
		unitPrice := order.NewMoney(itemReq.UnitPrice, itemReq.Currency)

		item, err := order.NewOrderItem(o.ID, itemReq.ProductID, itemReq.ProductName, itemReq.Quantity, unitPrice)
		if err != nil {
			return nil, err
		}
		item.ID = itemID

		if err := o.AddItem(item); err != nil {
			return nil, err
		}
	}

	// 保存订单
	if err := s.orderRepo.Create(ctx, o); err != nil {
		return nil, err
	}

	return domainOrderToDTO(o), nil
}

// CancelOrder 取消订单（命令）
func (s *Service) CancelOrder(ctx context.Context, orderID string) error {
	o, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if err := o.Cancel(); err != nil {
		return err
	}

	return s.orderRepo.Update(ctx, o)
}

// ProcessPayment 处理支付（命令）
func (s *Service) ProcessPayment(ctx context.Context, orderID string, req ProcessPaymentRequest) (*PaymentDTO, error) {
	// 验证订单是否可以支付
	if err := s.orderService.ValidateOrderForPayment(ctx, orderID); err != nil {
		return nil, err
	}

	// 获取订单
	o, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// 创建支付
	method := order.PaymentMethod(req.Method)
	payment, err := order.NewPayment(orderID, o.TotalAmount, method)
	if err != nil {
		return nil, err
	}

	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	payment.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	// 调用支付网关
	transactionID, gatewayResponse, err := s.paymentGateway.ProcessPayment(ctx, o.TotalAmount, method)
	if err != nil {
		if markErr := payment.MarkAsFailed(err.Error()); markErr != nil {
			return nil, markErr
		}
		if err := s.paymentRepo.Create(ctx, payment); err != nil {
			return nil, err
		}
		return nil, order.ErrPaymentFailed
	}

	// 标记支付成功
	if err := payment.MarkAsCompleted(transactionID, gatewayResponse); err != nil {
		return nil, err
	}

	// 保存支付
	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	// 更新订单状态为已支付
	if err := o.MarkAsPaid(); err != nil {
		return nil, err
	}

	if err := s.orderRepo.Update(ctx, o); err != nil {
		return nil, err
	}

	return domainPaymentToDTO(payment), nil
}

// RefundPayment 退款（命令）
func (s *Service) RefundPayment(ctx context.Context, orderID string) error {
	// 验证是否可以退款
	if err := s.orderService.ValidateRefund(ctx, orderID); err != nil {
		return err
	}

	// 获取支付记录
	payment, err := s.paymentRepo.FindByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	// 调用支付网关退款
	if err := s.paymentGateway.RefundPayment(ctx, payment.TransactionID, payment.Amount); err != nil {
		return err
	}

	// 更新支付状态
	if err := payment.Refund(); err != nil {
		return err
	}

	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return err
	}

	// 更新订单状态
	o, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}

	if err := o.Refund(); err != nil {
		return err
	}

	return s.orderRepo.Update(ctx, o)
}

// CreateShipment 创建发货（命令）
func (s *Service) CreateShipment(ctx context.Context, orderID string, req CreateShipmentRequest) (*ShipmentDTO, error) {
	// 验证订单是否可以发货
	if err := s.orderService.ValidateOrderForShipment(ctx, orderID); err != nil {
		return nil, err
	}

	// 创建地址
	address, err := order.NewAddress(
		req.Address.Street,
		req.Address.City,
		req.Address.State,
		req.Address.PostalCode,
		req.Address.Country,
	)
	if err != nil {
		return nil, err
	}

	// 创建发货
	shipment, err := order.NewShipment(orderID, address, req.ShippingMethod)
	if err != nil {
		return nil, err
	}

	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	shipment.ID = ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	// 计算预计送达时间
	estimatedDate := s.orderService.CalculateEstimatedDelivery(req.ShippingMethod)
	shipment.SetEstimatedDeliveryDate(estimatedDate)

	// 保存发货
	if err := s.shipmentRepo.Create(ctx, shipment); err != nil {
		return nil, err
	}

	return domainShipmentToDTO(shipment), nil
}

// UpdateShipment 更新发货（命令）
func (s *Service) UpdateShipment(ctx context.Context, shipmentID string, req UpdateShipmentRequest) (*ShipmentDTO, error) {
	shipment, err := s.shipmentRepo.FindByID(ctx, shipmentID)
	if err != nil {
		return nil, err
	}

	// 开始处理
	if err := shipment.StartProcessing(); err != nil {
		return nil, err
	}

	// 发货
	if err := shipment.Ship(req.TrackingNumber, req.Carrier); err != nil {
		return nil, err
	}

	// 保存更新
	if err := s.shipmentRepo.Update(ctx, shipment); err != nil {
		return nil, err
	}

	// 更新订单状态为已完成
	o, err := s.orderRepo.FindByID(ctx, shipment.OrderID)
	if err != nil {
		return nil, err
	}

	if err := o.Complete(); err != nil {
		return nil, err
	}

	if err := s.orderRepo.Update(ctx, o); err != nil {
		return nil, err
	}

	return domainShipmentToDTO(shipment), nil
}
