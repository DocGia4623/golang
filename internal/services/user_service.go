package services

import "testwire/internal/models"

type UserSerive interface {
	AddRole(uint, []models.Role) error
	FindRole([]string) ([]models.Role, error)
	FindAll() ([]models.User, error)
}
