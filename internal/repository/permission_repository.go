package repository

import (
	"testwire/internal/models"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	FindIfExist(string) (*models.Permission, error)
	// role func
	FindRoleBelong(string) ([]models.Role, error)
}

type PermissionRepositoryImpl struct {
	Db *gorm.DB
}

func NewPermissionRepositoryImpl(db *gorm.DB) PermissionRepository {
	return &PermissionRepositoryImpl{Db: db}
}

func (p *PermissionRepositoryImpl) FindIfExist(permissionName string) (*models.Permission, error) {
	var permission models.Permission
	err := p.Db.Where("name = ?", permissionName).First(&permission).Error

	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (p *PermissionRepositoryImpl) FindRoleBelong(permission string) ([]models.Role, error) {
	var roles []models.Role
	result := p.Db.Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("permissions.name = ?", permission).
		Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (p *PermissionRepositoryImpl) FindRole(roles []string) ([]models.Role, error) {
	var foundRoles []models.Role

	if err := p.Db.Where("name IN ?", roles).Find(&foundRoles).Error; err != nil {
		return nil, err
	}

	return foundRoles, nil
}
