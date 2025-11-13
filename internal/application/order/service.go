package order

import (
	"context"

	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/order"
)

// PaymentGateway 支付网关接口（端口）
type PaymentGateway interface {
	ProcessPayment(ctx context.Context, amount order.Money, method order.PaymentMethod) (transactionID, response string, err error)
	RefundPayment(ctx context.Context, transactionID string, amount order.Money) error
}

// Service 订单应用服务
type Service struct {
	orderRepo    order.OrderRepository
	paymentRepo  order.PaymentRepository
	shipmentRepo order.ShipmentRepository
	orderService *order.Service

	paymentGateway PaymentGateway
}

// NewService 创建订单应用服务
func NewService(
	orderRepo order.OrderRepository,
	paymentRepo order.PaymentRepository,
	shipmentRepo order.ShipmentRepository,
	orderService *order.Service,
	paymentGateway PaymentGateway,
) *Service {
	return &Service{
		orderRepo:      orderRepo,
		paymentRepo:    paymentRepo,
		shipmentRepo:   shipmentRepo,
		orderService:   orderService,
		paymentGateway: paymentGateway,
	}
}

// domainOrderToDTO 转换订单为DTO
func domainOrderToDTO(o *order.Order) *OrderDTO {
	items := make([]*OrderItemDTO, len(o.Items))
	for i, item := range o.Items {
		items[i] = &OrderItemDTO{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice: MoneyDTO{
				Amount:   item.UnitPrice.Amount,
				Currency: item.UnitPrice.Currency,
			},
			Subtotal: MoneyDTO{
				Amount:   item.Subtotal.Amount,
				Currency: item.Subtotal.Currency,
			},
		}
	}

	return &OrderDTO{
		ID:          o.ID,
		UserID:      o.UserID,
		OrderNumber: o.OrderNumber,
		Status:      string(o.Status),
		Items:       items,
		TotalAmount: MoneyDTO{
			Amount:   o.TotalAmount.Amount,
			Currency: o.TotalAmount.Currency,
		},
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

// domainPaymentToDTO 转换支付为DTO
func domainPaymentToDTO(p *order.Payment) *PaymentDTO {
	return &PaymentDTO{
		ID:      p.ID,
		OrderID: p.OrderID,
		Amount: MoneyDTO{
			Amount:   p.Amount.Amount,
			Currency: p.Amount.Currency,
		},
		Method:          string(p.Method),
		Status:          string(p.Status),
		TransactionID:   p.TransactionID,
		GatewayResponse: p.GatewayResponse,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

// domainShipmentToDTO 转换发货为DTO
func domainShipmentToDTO(s *order.Shipment) *ShipmentDTO {
	return &ShipmentDTO{
		ID:             s.ID,
		OrderID:        s.OrderID,
		TrackingNumber: s.TrackingNumber,
		Carrier:        s.Carrier,
		ShippingMethod: s.ShippingMethod,
		Address: AddressDTO{
			Street:     s.Address.Street,
			City:       s.Address.City,
			State:      s.Address.State,
			PostalCode: s.Address.PostalCode,
			Country:    s.Address.Country,
		},
		Status:        string(s.Status),
		EstimatedDate: s.EstimatedDate,
		ShippedAt:     s.ShippedAt,
		DeliveredAt:   s.DeliveredAt,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}
