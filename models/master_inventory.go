package model

import "gorm.io/gorm"

type Inventory struct {
	gorm.Model
	Name          string `json:"name" validate:"required"`
	DistributorID uint   `json:"distributor_id" validate:"required"`
	ProductID     uint   `json:"product_id" validate:"required"`
	Distributor   Distributor
	Product       Product
	Stock         uint    `json:"stock" validate:"required"`
	Price         float64 `json:"price" validate:"required"`
}

type InventoryResponse struct {
	ID                  uint                `json:"id"`
	Name                string              `json:"name"`
	DistributorID       uint                `json:"distributor_id"`
	ProductID           uint                `json:"product_id"`
	DistributorResponse DistributorResponse `json:"distributor"`
	ProductResponse     ProductResponse     `json:"product"`
	Stock               uint                `json:"stock"`
	Price               float64             `json:"price"`
	CreatedAt           string              `json:"created_at"`
	UpdatedAt           string              `json:"updated_at"`
}

type ErrValidationResp struct {
	Name string `json:"name"`
}

func (Inventory) TableName() string {
	return "master_inventories"
}
