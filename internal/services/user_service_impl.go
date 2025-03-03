package services

import (
	"testwire/internal/models"
	"testwire/internal/repository"
)

type UserServiceImpl struct {
	UserRepo repository.UserRepository
}

func NewUserServiceImpl(userRepo repository.UserRepository) UserSerive {
	return &UserServiceImpl{UserRepo: userRepo}
}

func (u *UserServiceImpl) AddRole(userId uint, roles []models.Role) error {
	return u.UserRepo.AddRole(userId, roles)
}

func (u *UserServiceImpl) FindRole(roles []string) ([]models.Role, error) {
	return u.UserRepo.FindRole(roles)
}

func (u *UserServiceImpl) FindAll() ([]models.User, error) {
	return u.UserRepo.FindAll()
}
