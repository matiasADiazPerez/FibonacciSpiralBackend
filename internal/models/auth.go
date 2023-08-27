package models

type AuthUser struct {
	Id       int    `json:"id" validate:"required"`
	Password string `json:"password" validate:"required"`
}
