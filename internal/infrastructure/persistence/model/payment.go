package model

import "time"

// Payment GORM支付模型
type Payment struct {
	ID              string    `gorm:"primaryKey;type:varchar(26)"`
	OrderID         string    `gorm:"uniqueIndex;not null;type:varchar(26)"`
	Amount          float64   `gorm:"not null;type:decimal(10,2)"`
	Currency        string    `gorm:"not null;type:varchar(3);default:'USD'"`
	Method          string    `gorm:"not null;type:varchar(20)"`
	Status          string    `gorm:"not null;type:varchar(20);default:'pending'"`
	TransactionID   string    `gorm:"type:varchar(255)"`
	GatewayResponse string    `gorm:"type:text"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`

	Order Order `gorm:"foreignKey:OrderID"`
}

// TableName 指定表名
func (Payment) TableName() string {
	return "payments"
}
