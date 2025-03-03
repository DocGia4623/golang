package models

import (
	"testwire/internal/constant"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID     uint                 `json:"user_id" gorm:"not null"`
	TotalPrice float64              `json:"total_price" gorm:"type:decimal(10,2);not null"`
	Status     constant.OrderStatus `json:"status" gorm:"default:0"`
	OrderItems []OrderDetail        `json:"order_items" gorm:"foreignKey:OrderID"`
}

type OrderDetail struct {
	gorm.Model
	OrderID   uint    `json:"order_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Quantity  uint    `json:"quantity" gorm:"not null;check:quantity>0"`
	Price     float64 `json:"price" gorm:"type:decimal(10,2);not null"`
}
