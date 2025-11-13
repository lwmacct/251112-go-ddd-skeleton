package model

import "time"

// OrderItem GORM订单项模型
type OrderItem struct {
	ID          string    `gorm:"primaryKey;type:varchar(26)"`
	OrderID     string    `gorm:"index;not null;type:varchar(26)"`
	ProductID   string    `gorm:"not null;type:varchar(26)"`
	ProductName string    `gorm:"not null;type:varchar(255)"`
	Quantity    int       `gorm:"not null;default:1"`
	UnitPrice   float64   `gorm:"not null;type:decimal(10,2)"`
	Subtotal    float64   `gorm:"not null;type:decimal(10,2)"`
	Currency    string    `gorm:"not null;type:varchar(3);default:'USD'"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	Order Order `gorm:"foreignKey:OrderID"`
}

// TableName 指定表名
func (OrderItem) TableName() string {
	return "order_items"
}
