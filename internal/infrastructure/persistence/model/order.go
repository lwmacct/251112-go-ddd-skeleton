package model

import "time"

// Order GORM订单模型
type Order struct {
	ID          string    `gorm:"primaryKey;type:varchar(26)"`
	UserID      string    `gorm:"index;not null;type:varchar(26)"`
	OrderNumber string    `gorm:"uniqueIndex;not null;type:varchar(50)"`
	Status      string    `gorm:"not null;type:varchar(20);default:'pending'"`
	TotalAmount float64   `gorm:"not null;type:decimal(10,2);default:0"`
	Currency    string    `gorm:"not null;type:varchar(3);default:'USD'"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// 关联关系
	User      User        `gorm:"foreignKey:UserID"`
	Items     []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	Payment   *Payment    `gorm:"foreignKey:OrderID"`
	Shipment  *Shipment   `gorm:"foreignKey:OrderID"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}
