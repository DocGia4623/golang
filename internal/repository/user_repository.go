package repository

import (
	"testwire/internal/dto/request"
	"testwire/internal/models"
)

type UserRepository interface {
	Save(request.CreateUserRequest) error
	Update(request.UpdateUserRequest, int) error
	Delete(int) error
	FindAll() ([]models.User, error)
	FindById(int) (models.User, error)
	FindByUsername(string) (*models.User, error)
	FindIfUserHasRole(uint, []models.Role) error
	AddRole(uint, []models.Role) error
	FindRole([]string) ([]models.Role, error)
}
