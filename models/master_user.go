package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
	RoleID   uint   `json:"role_id" validate:"required"`
	Role     Role
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID        uint         `json:"id"`
	Name      string       `json:"name"`
	RoleID    uint         `json:"role_id"`
	Username  string       `json:"username"`
	Role      RoleResponse `json:"role"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
}

type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (User) TableName() string {
	return "master_users"
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (user *User) CheckPasswordHash(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return err
	}

	return nil
}
