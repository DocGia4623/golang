package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint   `json:"ID"`
	UserName string `gorm:"unique" json:"username"`
	FullName string `json:"fullname"`
	Password string `json:"password"`
	Age      uint   `json:"age"`
	Roles    []Role `gorm:"many2many:user_roles"`
}

// Role struct
type Role struct {
	gorm.Model
	Name        string       `gorm:"unique"`
	Permissions []Permission `gorm:"many2many:role_permissions"`
}

// Permission struct
type Permission struct {
	gorm.Model
	Name string `gorm:"unique"`
}
