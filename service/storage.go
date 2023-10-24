package service

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func DBconnection() error {
	err := godotenv.Load()
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
		return err
	}

	sqlDB, _ := db.DB()
	err = sqlDB.Close()
	if err != nil {
		return err
	}

	return nil
}
