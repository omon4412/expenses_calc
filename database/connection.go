package database

import (
	"fmt"
	"project/config"
	"project/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	getConfig := config.GetConfig()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s port=%d sslmode=disable",
		getConfig.Database.Host, getConfig.Database.User, getConfig.Database.Password, getConfig.Database.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB = db

	db.AutoMigrate(&models.User{}, &models.Category{}, &models.Expense{})
	return db, nil
}
