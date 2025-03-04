package services

import (
	"testwire/internal/models"
	"testwire/internal/repository"
)

type OrderSerivce interface {
	CreateOrder(models.Order) error
	CancelOrder(uint) error
	FindOrderByUserID(uint) ([]models.Order, error)
}

type OrderSerivceImpl struct {
	OrderRepo repository.OrderRepository
}

func NewOrderServiceImpl(orderRepo repository.OrderRepository) OrderSerivce {
	return &OrderSerivceImpl{OrderRepo: orderRepo}
}

func (o *OrderSerivceImpl) CreateOrder(order models.Order) error {
	return o.OrderRepo.CreateOrder(order)
}
func (o *OrderSerivceImpl) CancelOrder(id uint) error {
	return o.OrderRepo.CancelOrder(id)
}
func (o *OrderSerivceImpl) FindOrderByUserID(id uint) ([]models.Order, error) {
	return o.OrderRepo.FindOrderByUserID(id)
}
