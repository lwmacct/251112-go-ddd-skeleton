package repository

import (
	"context"
	"errors"

	"github.com/lwmacct/251112-go-ddd-skeleton/internal/domain/order"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/mapper"
	"github.com/lwmacct/251112-go-ddd-skeleton/internal/infrastructure/persistence/model"
	"gorm.io/gorm"
)

// OrderRepository 订单仓储实现
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository 创建订单仓储
func NewOrderRepository(db *gorm.DB) order.OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, o *order.Order) error {
	m := mapper.OrderToModel(o)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *OrderRepository) Update(ctx context.Context, o *order.Order) error {
	m := mapper.OrderToModel(o)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *OrderRepository) FindByID(ctx context.Context, id string) (*order.Order, error) {
	var m model.Order
	if err := r.db.WithContext(ctx).Preload("Items").First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, order.ErrOrderNotFound
		}
		return nil, err
	}
	return mapper.OrderToDomain(&m), nil
}

func (r *OrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*order.Order, error) {
	var m model.Order
	if err := r.db.WithContext(ctx).Preload("Items").First(&m, "order_number = ?", orderNumber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, order.ErrOrderNotFound
		}
		return nil, err
	}
	return mapper.OrderToDomain(&m), nil
}

func (r *OrderRepository) ListByUserID(ctx context.Context, userID string, offset, limit int) ([]*order.Order, int64, error) {
	var models []model.Order
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	orders := make([]*order.Order, len(models))
	for i, m := range models {
		orders[i] = mapper.OrderToDomain(&m)
	}

	return orders, total, nil
}

func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Order{}, "id = ?", id).Error
}

// PaymentRepository 支付仓储实现
type PaymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository 创建支付仓储
func NewPaymentRepository(db *gorm.DB) order.PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *order.Payment) error {
	m := mapper.PaymentToModel(payment)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PaymentRepository) Update(ctx context.Context, payment *order.Payment) error {
	m := mapper.PaymentToModel(payment)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *PaymentRepository) FindByID(ctx context.Context, id string) (*order.Payment, error) {
	var m model.Payment
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, order.ErrPaymentNotFound
		}
		return nil, err
	}
	return mapper.PaymentToDomain(&m), nil
}

func (r *PaymentRepository) FindByOrderID(ctx context.Context, orderID string) (*order.Payment, error) {
	var m model.Payment
	if err := r.db.WithContext(ctx).First(&m, "order_id = ?", orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, order.ErrPaymentNotFound
		}
		return nil, err
	}
	return mapper.PaymentToDomain(&m), nil
}

func (r *PaymentRepository) FindByTransactionID(ctx context.Context, transactionID string) (*order.Payment, error) {
	var m model.Payment
	if err := r.db.WithContext(ctx).First(&m, "transaction_id = ?", transactionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, order.ErrPaymentNotFound
		}
		return nil, err
	}
	return mapper.PaymentToDomain(&m), nil
}

// ShipmentRepository 发货仓储实现
type ShipmentRepository struct {
	db *gorm.DB
}

// NewShipmentRepository 创建发货仓储
func NewShipmentRepository(db *gorm.DB) order.ShipmentRepository {
	return &ShipmentRepository{db: db}
}

func (r *ShipmentRepository) Create(ctx context.Context, shipment *order.Shipment) error {
	m := mapper.ShipmentToModel(shipment)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *ShipmentRepository) Update(ctx context.Context, shipment *order.Shipment) error {
	m := mapper.ShipmentToModel(shipment)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *ShipmentRepository) FindByID(ctx context.Context, id string) (*order.Shipment, error) {
	var m model.Shipment
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, order.ErrShipmentNotFound
		}
		return nil, err
	}
	return mapper.ShipmentToDomain(&m)
}

func (r *ShipmentRepository) FindByOrderID(ctx context.Context, orderID string) (*order.Shipment, error) {
	var m model.Shipment
	if err := r.db.WithContext(ctx).First(&m, "order_id = ?", orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, order.ErrShipmentNotFound
		}
		return nil, err
	}
	return mapper.ShipmentToDomain(&m)
}

func (r *ShipmentRepository) FindByTrackingNumber(ctx context.Context, trackingNumber string) (*order.Shipment, error) {
	var m model.Shipment
	if err := r.db.WithContext(ctx).First(&m, "tracking_number = ?", trackingNumber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, order.ErrShipmentNotFound
		}
		return nil, err
	}
	return mapper.ShipmentToDomain(&m)
}
