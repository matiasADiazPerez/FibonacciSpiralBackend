package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Password string `json:",omitempty"`
}

type CreateUser struct {
	Name     string
	Password string
}

type ChangePassword struct {
	CurrentPassword string
	NewPassword     string
}
