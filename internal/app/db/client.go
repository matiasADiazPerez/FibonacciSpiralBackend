package db

import (
	"fmt"
	"os"
	"spiralmatrix/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newDB() (*gorm.DB, error) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	if user == "" || password == "" || host == "" {
		return nil, fmt.Errorf("db env variables are not set")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s ", host, user, password)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func migrateDB(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{})
	fmt.Println(db.Migrator().HasTable(&models.User{}))
	return err
}

func InitDb() (*gorm.DB, error) {
	db, err := newDB()
	if err != nil {
		return nil, err
	}
	return db, migrateDB(db)
}
