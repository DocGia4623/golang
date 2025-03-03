package repository

import "testwire/internal/models"

type PermissionRepository interface {
	FindIfExist(string) (*models.Permission, error)
	// role func
	FindRoleBelong(string) ([]models.Role, error)
}
