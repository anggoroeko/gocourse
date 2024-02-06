package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name string `json:"name"`
}

type ProductResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Product) TableName() string {
	return "master_products"
}
