package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Password string `json:",omitempty"`
}

type CreateUser struct {
	Name     string `validate:"required"`
	Password string `validate:"required"`
}

type ChangePassword struct {
	CurrentPassword string `validate:"required"`
	NewPassword     string `validate:"required"`
}
