package model

import "gorm.io/gorm"

type Purchase struct {
	gorm.Model
	InventoryID uint      `json:"inventory_id"`
	Date        string    `json:"date"`
	Quantity    uint      `json:"quantity"`
	Inventory   Inventory `validate:"-"`
	Description string    `json:"description"`
}

type PurchaseResponse struct {
	ID          uint              `json:"id"`
	InventoryID uint              `json:"inventory_id"`
	Date        string            `json:"date"`
	Quantity    uint              `json:"quantity"`
	Inventory   InventoryResponse `json:"inventory"`
	Description string            `json:"description"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}
