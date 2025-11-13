package model

import "time"

// Shipment GORM发货模型
type Shipment struct {
	ID             string     `gorm:"primaryKey;type:varchar(26)"`
	OrderID        string     `gorm:"uniqueIndex;not null;type:varchar(26)"`
	TrackingNumber string     `gorm:"type:varchar(100)"`
	Carrier        string     `gorm:"type:varchar(100)"`
	ShippingMethod string     `gorm:"not null;type:varchar(50)"`
	Street         string     `gorm:"not null;type:varchar(255)"`
	City           string     `gorm:"not null;type:varchar(100)"`
	State          string     `gorm:"type:varchar(100)"`
	PostalCode     string     `gorm:"type:varchar(20)"`
	Country        string     `gorm:"not null;type:varchar(100)"`
	Status         string     `gorm:"not null;type:varchar(20);default:'pending'"`
	EstimatedDate  *time.Time
	ShippedAt      *time.Time
	DeliveredAt    *time.Time
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`

	Order Order `gorm:"foreignKey:OrderID"`
}

// TableName 指定表名
func (Shipment) TableName() string {
	return "shipments"
}
