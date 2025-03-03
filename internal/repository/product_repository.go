package repository

import "testwire/internal/models"

type ProductRepository interface {
	Save(models.Product) error
	Find(string) (*models.Product, error)
	Delete(string) error
}
