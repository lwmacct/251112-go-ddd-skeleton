package mapper

import (
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/order"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/model"
)

// OrderToModel 转换订单到模型
func OrderToModel(o *order.Order) *model.Order {
	m := &model.Order{
		ID:          o.ID,
		UserID:      o.UserID,
		OrderNumber: o.OrderNumber,
		Status:      string(o.Status),
		TotalAmount: o.TotalAmount.Amount,
		Currency:    o.TotalAmount.Currency,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}

	// 转换订单项
	items := make([]model.OrderItem, len(o.Items))
	for i, item := range o.Items {
		items[i] = model.OrderItem{
			ID:          item.ID,
			OrderID:     item.OrderID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice.Amount,
			Subtotal:    item.Subtotal.Amount,
			Currency:    item.UnitPrice.Currency,
			CreatedAt:   item.CreatedAt,
		}
	}
	m.Items = items

	return m
}

// OrderToDomain 转换模型到订单
func OrderToDomain(m *model.Order) *order.Order {
	o := &order.Order{
		ID:          m.ID,
		UserID:      m.UserID,
		OrderNumber: m.OrderNumber,
		Status:      order.OrderStatus(m.Status),
		TotalAmount: order.NewMoney(m.TotalAmount, m.Currency),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}

	// 转换订单项
	items := make([]*order.OrderItem, len(m.Items))
	for i, item := range m.Items {
		items[i] = &order.OrderItem{
			ID:          item.ID,
			OrderID:     item.OrderID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitPrice:   order.NewMoney(item.UnitPrice, item.Currency),
			Subtotal:    order.NewMoney(item.Subtotal, item.Currency),
			CreatedAt:   item.CreatedAt,
		}
	}
	o.Items = items

	return o
}

// PaymentToModel 转换支付到模型
func PaymentToModel(p *order.Payment) *model.Payment {
	return &model.Payment{
		ID:              p.ID,
		OrderID:         p.OrderID,
		Amount:          p.Amount.Amount,
		Currency:        p.Amount.Currency,
		Method:          string(p.Method),
		Status:          string(p.Status),
		TransactionID:   p.TransactionID,
		GatewayResponse: p.GatewayResponse,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

// PaymentToDomain 转换模型到支付
func PaymentToDomain(m *model.Payment) *order.Payment {
	return &order.Payment{
		ID:              m.ID,
		OrderID:         m.OrderID,
		Amount:          order.NewMoney(m.Amount, m.Currency),
		Method:          order.PaymentMethod(m.Method),
		Status:          order.PaymentStatus(m.Status),
		TransactionID:   m.TransactionID,
		GatewayResponse: m.GatewayResponse,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

// ShipmentToModel 转换发货到模型
func ShipmentToModel(s *order.Shipment) *model.Shipment {
	return &model.Shipment{
		ID:             s.ID,
		OrderID:        s.OrderID,
		TrackingNumber: s.TrackingNumber,
		Carrier:        s.Carrier,
		ShippingMethod: s.ShippingMethod,
		Street:         s.Address.Street,
		City:           s.Address.City,
		State:          s.Address.State,
		PostalCode:     s.Address.PostalCode,
		Country:        s.Address.Country,
		Status:         string(s.Status),
		EstimatedDate:  s.EstimatedDate,
		ShippedAt:      s.ShippedAt,
		DeliveredAt:    s.DeliveredAt,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
}

// ShipmentToDomain 转换模型到发货
func ShipmentToDomain(m *model.Shipment) (*order.Shipment, error) {
	address, err := order.NewAddress(m.Street, m.City, m.State, m.PostalCode, m.Country)
	if err != nil {
		return nil, err
	}

	return &order.Shipment{
		ID:             m.ID,
		OrderID:        m.OrderID,
		TrackingNumber: m.TrackingNumber,
		Carrier:        m.Carrier,
		ShippingMethod: m.ShippingMethod,
		Address:        address,
		Status:         order.ShipmentStatus(m.Status),
		EstimatedDate:  m.EstimatedDate,
		ShippedAt:      m.ShippedAt,
		DeliveredAt:    m.DeliveredAt,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}, nil
}
