package model

import "gorm.io/gorm"

type Sales struct {
	gorm.Model
	InventoryID uint      `json:"inventory_id" validate:"required"`
	MemberID    uint      `json:"member_id" validate:"required"`
	Quantity    uint      `json:"quantity" validate:"required"`
	ChasierID   uint      `json:"chasier_id" validate:"required"`
	Inventory   Inventory `validate:"-"`
	Member      Member
	Chasier     Chasier
}

type SalesResponse struct {
	ID          uint   `json:"id"`
	Member      string `json:"member"`
	Chasier     string `json:"chasier"`
	Inventory   string `json:"inventory"`
	Distributor string `json:"distributor"`
	Product     string `json:"product"`
	Quantity    uint   `json:"quantity"`
	Price       uint   `json:"price"`
}
