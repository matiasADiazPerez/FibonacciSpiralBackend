package main

import (
	"log"
	_ "spiralmatrix/docs"
	"spiralmatrix/internal/app/auth"
	"spiralmatrix/internal/app/db"
	"spiralmatrix/internal/app/server"
	"spiralmatrix/internal/app/spiral"
	"spiralmatrix/internal/app/user"

	"github.com/joho/godotenv"
)

//	@title			Spiral Matrix API
//	@version		1.0
//	@description	The backend of the fibonacci spiral matrix implementation
//	@BasePath		/
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
	spiralHandler := spiral.NewSpiralHandler()
	server.Start(&userHandler, &authHandler, &spiralHandler)
}
