package repository

import (
	"errors"
	"testwire/internal/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAll() ([]models.Product, error)
	Save(models.Product) error
	Find(string) (*models.Product, error)
	Delete(string) error
}

type ProductRepositoryImpl struct {
	Db *gorm.DB
}

func NewProductRepositoryImpl(db *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{Db: db}
}

func (p *ProductRepositoryImpl) Save(product models.Product) error {
	result := p.Db.Save(&product)
	return result.Error
}
func (p *ProductRepositoryImpl) Find(name string) (*models.Product, error) {
	var product models.Product
	result := p.Db.Where("name = ?", name).First(&product)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Trả về nil thay vì lỗi nếu không tìm thấy
	}
	return &product, result.Error
}
func (p *ProductRepositoryImpl) Delete(name string) error {
	result := p.Db.Where("name = ?", name).Delete(&models.Product{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Nếu không xóa được thì trả về lỗi
	}
	return result.Error
}

func (p *ProductRepositoryImpl) GetAll() ([]models.Product, error) {
	var products []models.Product
	result := p.Db.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}
