package model

import (
	"database/sql"

	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	Name      string       `json:"name"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

type MemberResponse struct {
	ID        uint         `json:"id"`
	Name      string       `json:"name"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

func (Member) TableName() string {
	return "master_members"
}
