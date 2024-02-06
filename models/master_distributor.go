package model

import (
	"time"

	"gorm.io/gorm"
)

type Distributor struct {
	gorm.Model
	Name string `json:"name"`
}

type DistributorResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Distributor) TableName() string {
	return "master_distributors"
}
