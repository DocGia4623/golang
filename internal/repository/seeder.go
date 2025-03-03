package repository

import (
	"log"
	"testwire/config"
	"testwire/internal/constant"
	"testwire/internal/models"
)

// Seed roles & permissions
func SeedRolesAndPermissions() {
	db := config.DB

	// Define permissions
	permissions := []models.Permission{
		{Name: constant.ReadUser},
		{Name: constant.WriteUser},
		{Name: constant.DeleteUser},
	}

	// Define roles
	roles := []models.Role{
		{Name: "admin"},
		{Name: "editor"},
		{Name: "viewer"},
	}

	// Insert into database
	db.Create(&permissions)
	db.Create(&roles)

	var dbPermissions []models.Permission
	db.Where("name IN ?", []string{
		constant.ReadUser,
		constant.WriteUser,
		constant.DeleteUser,
	}).Find(&dbPermissions)

	// Assign permissions to roles (Admin gets all)
	var adminRole models.Role
	if err := db.First(&adminRole, "name = ?", "admin").Error; err == nil {
		// Kiểm tra nếu chưa có quan hệ thì mới Append
		if err := db.Model(&adminRole).Association("Permissions").Error; err == nil {
			db.Model(&adminRole).Association("Permissions").Append(dbPermissions)
		}
	}

	log.Println("✅ Seeded roles and permissions")
}
