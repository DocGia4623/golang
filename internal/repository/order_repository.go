package repository

import "testwire/internal/models"

type OrderRepository interface {
	CreateOrder(models.Order) error
	CancelOrder(uint) error
	FindOrderByUserID(uint) ([]models.Order, error)
}
