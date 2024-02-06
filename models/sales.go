package model

import "gorm.io/gorm"

type Sales struct {
	gorm.Model
	InventoryID uint `json:"inventory_id"`
	MemberID    uint `json:"member_id"`
	Quantity    uint `json:"quantity"`
	ChasierID   uint `json:"cashier_id"`
	Inventory   Inventory
	Member      Member
	Chasier     Chasier
}

type SalesResponse struct {
	InventoryID uint `json:"inventory_id"`
	MemberID    uint `json:"member_id"`
	Quantity    uint `json:"quantity"`
	ChasierID   uint `json:"cashier_id"`
	Inventory   Inventory
	Member      Member
	Chasier     Chasier
}
