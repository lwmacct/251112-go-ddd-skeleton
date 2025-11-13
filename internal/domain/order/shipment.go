package order

import (
	"errors"
	"time"
)

// ShipmentStatus 发货状态
type ShipmentStatus string

const (
	ShipmentStatusPending    ShipmentStatus = "pending"
	ShipmentStatusProcessing ShipmentStatus = "processing"
	ShipmentStatusShipped    ShipmentStatus = "shipped"
	ShipmentStatusDelivered  ShipmentStatus = "delivered"
	ShipmentStatusCancelled  ShipmentStatus = "cancelled"
)

// Address 地址值对象
type Address struct {
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}

// NewAddress 创建地址
func NewAddress(street, city, state, postalCode, country string) (Address, error) {
	if street == "" || city == "" || country == "" {
		return Address{}, errors.New("street, city, and country are required")
	}
	return Address{
		Street:     street,
		City:       city,
		State:      state,
		PostalCode: postalCode,
		Country:    country,
	}, nil
}

// FullAddress 返回完整地址字符串
func (a Address) FullAddress() string {
	return a.Street + ", " + a.City + ", " + a.State + " " + a.PostalCode + ", " + a.Country
}

// Shipment 发货实体
type Shipment struct {
	ID             string
	OrderID        string
	TrackingNumber string
	Carrier        string
	ShippingMethod string
	Address        Address
	Status         ShipmentStatus
	EstimatedDate  *time.Time
	ShippedAt      *time.Time
	DeliveredAt    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewShipment 创建发货
func NewShipment(orderID string, address Address, shippingMethod string) (*Shipment, error) {
	if orderID == "" {
		return nil, errors.New("orderID cannot be empty")
	}

	return &Shipment{
		OrderID:        orderID,
		Address:        address,
		ShippingMethod: shippingMethod,
		Status:         ShipmentStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

// StartProcessing 开始处理
func (s *Shipment) StartProcessing() error {
	if s.Status != ShipmentStatusPending {
		return errors.New("can only start processing pending shipments")
	}
	s.Status = ShipmentStatusProcessing
	s.UpdatedAt = time.Now()
	return nil
}

// Ship 发货
func (s *Shipment) Ship(trackingNumber, carrier string) error {
	if s.Status != ShipmentStatusProcessing {
		return errors.New("can only ship processing shipments")
	}
	if trackingNumber == "" {
		return errors.New("tracking number is required")
	}

	s.Status = ShipmentStatusShipped
	s.TrackingNumber = trackingNumber
	s.Carrier = carrier
	now := time.Now()
	s.ShippedAt = &now
	s.UpdatedAt = now
	return nil
}

// Deliver 确认送达
func (s *Shipment) Deliver() error {
	if s.Status != ShipmentStatusShipped {
		return errors.New("can only deliver shipped shipments")
	}

	s.Status = ShipmentStatusDelivered
	now := time.Now()
	s.DeliveredAt = &now
	s.UpdatedAt = now
	return nil
}

// Cancel 取消发货
func (s *Shipment) Cancel() error {
	if s.Status == ShipmentStatusDelivered {
		return errors.New("cannot cancel delivered shipments")
	}
	s.Status = ShipmentStatusCancelled
	s.UpdatedAt = time.Now()
	return nil
}

// SetEstimatedDeliveryDate 设置预计送达时间
func (s *Shipment) SetEstimatedDeliveryDate(date time.Time) {
	s.EstimatedDate = &date
	s.UpdatedAt = time.Now()
}

// IsDelivered 是否已送达
func (s *Shipment) IsDelivered() bool {
	return s.Status == ShipmentStatusDelivered
}

// CanBeCancelled 是否可以取消
func (s *Shipment) CanBeCancelled() bool {
	return s.Status != ShipmentStatusDelivered
}

