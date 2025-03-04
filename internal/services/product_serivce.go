package services

import (
	"testwire/internal/models"
	"testwire/internal/repository"
)

type ProductService interface {
	SaveProduct(models.Product) error
	FindProduct(string) (*models.Product, error)
	DeleteProduct(string) error
	GetAll() ([]models.Product, error)
}

type ProductServiceImpl struct {
	ProductRepo repository.ProductRepository
}

func NewProductServiceImpl(productRepo repository.ProductRepository) ProductService {
	return &ProductServiceImpl{ProductRepo: productRepo}
}

func (p *ProductServiceImpl) SaveProduct(product models.Product) error {
	return p.ProductRepo.Save(product)
}
func (p *ProductServiceImpl) FindProduct(name string) (*models.Product, error) {
	return p.ProductRepo.Find(name)
}
func (p *ProductServiceImpl) DeleteProduct(name string) error {
	return p.ProductRepo.Delete(name)
}
func (p *ProductServiceImpl) GetAll() ([]models.Product, error) {
	return p.ProductRepo.GetAll()
}
