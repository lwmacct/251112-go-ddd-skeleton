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

