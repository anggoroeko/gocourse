package model

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name string `json:"name"`
}

type RoleResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (Role) TableName() string {
	return "master_roles"
}
