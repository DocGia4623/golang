package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ID    uint    `json:"ID"`
	Name  string  `json:"name" gorm:"not null"`
	Price float32 `json:"price" gorm:"not null;default:0"`
	Stock uint    `json:"stock" gorm:"not null;default:0"`
}
