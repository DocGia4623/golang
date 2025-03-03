package repository

import (
	"errors"
	"testwire/internal/constant"
	"testwire/internal/models"

	"gorm.io/gorm"
)

type OrderRepositoryImpl struct {
	Db *gorm.DB
}

func NewOrderRepositoryImpl(db *gorm.DB) OrderRepository {
	return &OrderRepositoryImpl{Db: db}
}

func (o *OrderRepositoryImpl) CreateOrder(order models.Order) error {
	result := o.Db.Create(order)
	return result.Error
}
func (o *OrderRepositoryImpl) CancelOrder(id uint) error {
	result := o.Db.Where("ID = ?", id).Update("status", constant.Cancelled)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
func (o *OrderRepositoryImpl) FindOrderByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	result := o.Db.Where("user_id = ?", userID).Find(&orders)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return orders, result.Error
}
