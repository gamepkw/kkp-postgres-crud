package model

import (
	"time"
)

type Order struct {
	ID           int         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CustomerName string      `gorm:"column:customer_name;type:varchar(100);not null" json:"customerName" validate:"required,max=100"`
	TotalAmount  float64     `gorm:"column:total_amount;type:decimal(10,2);not null" json:"totalAmount"`
	Status       string      `gorm:"column:status;type:varchar(20);not null" json:"status" validate:"required,max=20"`
	CreatedAt    time.Time   `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time   `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	OrderItems   []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items" validate:"required,min=1,dive"`
}

func (r *Order) TableName() string {
	return "orders.orders"
}

type OrderItem struct {
	ID          int     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OrderID     int     `gorm:"column:order_id;not null;index" json:"orderID"`
	ProductName string  `gorm:"column:product_name;type:varchar(100);not null" json:"productName" validate:"required,max=100"`
	Quantity    int     `gorm:"column:quantity;not null" json:"quantity" validate:"required,min=1"`
	Price       float64 `gorm:"column:price;type:decimal(10,2);not null" json:"price" validate:"required,min=0"`
}

func (r *OrderItem) TableName() string {
	return "orders.order_items"
}
