package model

import "time"

// Invoice GORM发票模型
type Invoice struct {
	ID          string    `gorm:"primaryKey;type:varchar(26)"`
	OrderID     string    `gorm:"uniqueIndex;not null;type:varchar(26)"`
	InvoiceNumber string  `gorm:"uniqueIndex;not null;type:varchar(50)"`
	Amount      float64   `gorm:"not null;type:decimal(10,2)"`
	Currency    string    `gorm:"not null;type:varchar(3);default:'USD'"`
	Status      string    `gorm:"not null;type:varchar(20);default:'draft'"`
	IssuedAt    *time.Time
	DueAt       *time.Time
	PaidAt      *time.Time
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Order Order `gorm:"foreignKey:OrderID"`
}

// TableName 指定表名
func (Invoice) TableName() string {
	return "invoices"
}
