package main

import (
	"log"
	"spiralmatrix/internal/app/auth"
	"spiralmatrix/internal/app/db"
	"spiralmatrix/internal/app/server"
	"spiralmatrix/internal/app/user"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	newDB, err := db.InitDb()
	if err != nil {
		panic(err)
	}
	userHandler := user.NewUserHandler(newDB)
	authHandler := auth.NewAuthHandler(newDB)
	server.Start(&userHandler, &authHandler)
}
