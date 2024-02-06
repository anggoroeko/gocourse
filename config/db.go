package config

import (
	model "go_pos_v1_2/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error

	dsn := "root:root@tcp(127.0.0.1:3306)/go_pos_v1_2_2023?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed To Connect Database")
	}

	DB.AutoMigrate(&model.Role{}, &model.User{}, &model.Chasier{}, &model.Member{}, &model.Product{}, &model.Distributor{}, &model.Inventory{}, &model.Sales{}, &model.Purchase{})
}
