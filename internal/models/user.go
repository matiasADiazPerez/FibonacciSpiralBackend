package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"size:200"`
	Email    string `json:"email" gorm:"unique; size:200"`
	Password string `json:"password,omitempty" gorm:"size:250"`
}

type CreateUser struct {
	Name     string `validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `validate:"required"`
}

type ChangePassword struct {
	CurrentPassword string `validate:"required"`
	NewPassword     string `validate:"required"`
}
