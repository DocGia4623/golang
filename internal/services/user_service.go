package services

import (
	"fmt"
	"strconv"
	"strings"
	"testwire/config"
	"testwire/internal/models"
	"testwire/internal/repository"
	"testwire/utils"
)

type UserSerive interface {
	AddRole(uint, []models.Role) error
	FindRole([]string) ([]models.Role, error)
	FindAll() ([]models.User, error)
	GetUserID(string) (int, error)
}

type UserServiceImpl struct {
	UserRepo repository.UserRepository
}

func NewUserServiceImpl(userRepo repository.UserRepository) UserSerive {
	return &UserServiceImpl{UserRepo: userRepo}
}

func (u *UserServiceImpl) GetUserID(token string) (int, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	config, err := config.LoadConfig()
	if err != nil {
		return 0, err
	}
	sub, err := utils.ValidateAccessToken(token, config.AccessTokenSecret)
	if err != nil {
		return 0, err
	}
	id, err_id := strconv.Atoi(fmt.Sprint(sub))
	if err_id != nil {
		return 0, err_id
	}
	return id, nil
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
