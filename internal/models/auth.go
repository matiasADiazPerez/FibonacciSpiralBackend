package models

type AuthUser struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
}
