package main

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"project/config"
	"project/database"
	"project/logging"
	"project/models"
	"project/routes"
	"time"
)

func main() {
	config.LoadConfig("config.yaml")
	dbconnect, err := database.ConnectDB()
	if err != nil {
		panic("could not connect to db")
	}
	logging.Logger.Info("Connection is successful")

	addDefaultCategories(dbconnect)

	port := config.GetConfig().Server.Port
	timeout := config.GetConfig().Server.Timeout
	app := fiber.New(fiber.Config{
		IdleTimeout: time.Duration(timeout) * time.Second,
	})

	routes.SetupRoutes(app)

	err = app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		logging.Logger.Fatal("Could not start server: ", zap.Error(err))
	}
}

func addDefaultCategories(db *gorm.DB) {
	for _, category := range models.DefaultCategories {
		var existingCategory models.Category
		err := db.Where("name = ?", category.Name).First(&existingCategory).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Ошибка при проверке категории:", err)
			continue
		}

		if existingCategory.ID == 0 {
			err := db.Create(&category).Error
			if err != nil {
				fmt.Println("Ошибка при добавлении категории:", err)
			} else {
				fmt.Println("Дефолтная категория добавлена:", category.Name)
			}
		}
	}
}
