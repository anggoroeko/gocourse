package model

import (
	"gorm.io/gorm"
)

type Chasier struct {
	gorm.Model
	Name   string `json:"name"`
	UserID uint   `json:"user_id"`
	User   User   `validate:"-"` //:: CAN ALSO USE STRUCTONLY
}

type ChasierResponse struct {
	ID        uint         `json:"id"`
	Name      string       `json:"name"`
	UserID    uint         `json:"user_id"`
	User      UserResponse `json:"user"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
}

func (Chasier) TableName() string {
	return "master_chasiers"
}
