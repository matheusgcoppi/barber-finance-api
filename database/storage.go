package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/matheusgcoppi/barber-finance-api/database/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type CustomDB struct {
	Db *gorm.DB
}

func NewPostgres() (*CustomDB, error) {
	err := godotenv.Load("/Users/matheusgcoppi/Development/Golang/barber-finance/.env")
	if err != nil {
		log.Fatal(".env file could not be loaded.")
	}
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf(
		`host=%s user=%s password=%s dbname=%s port=%v sslmode=%s`,
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return nil, err
	}

	return &CustomDB{Db: db}, nil
}
