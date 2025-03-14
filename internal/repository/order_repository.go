package repository

import (
	"errors"
	"testwire/internal/constant"
	"testwire/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(models.Order) error
	CancelOrder(uint) error
	FindOrderByUserID(uint) ([]models.Order, error)
}

type OrderRepositoryImpl struct {
	Db *gorm.DB
}

func NewOrderRepositoryImpl(db *gorm.DB) OrderRepository {
	return &OrderRepositoryImpl{Db: db}
}

func (o *OrderRepositoryImpl) CreateOrder(order models.Order) error {
	// Bắt đầu transaction để đảm bảo tính toàn vẹn của dữ liệu
	tx := o.Db.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback() // Nếu gặp lỗi, rollback transaction
		return err
	}

	// Lưu các chi tiết đơn hàng
	for _, item := range order.OrderItems {
		item.OrderID = order.ID // Gán OrderID cho OrderDetail
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback() // Nếu có lỗi khi lưu chi tiết đơn hàng, rollback
			return err
		}
	}

	// Commit transaction nếu mọi thứ ổn
	return tx.Commit().Error
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
