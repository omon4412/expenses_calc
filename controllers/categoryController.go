package controllers

import (
	"github.com/gofiber/fiber/v3"
	"project/database"
	"project/logging"
	"project/models"
)

func GetCategories(c fiber.Ctx) error {
	logging.Logger.Info("Request to get categories")

	id, err2, done := CheckUser(c)
	if done {
		return err2
	}
	var categories []models.Category
	database.DB.Where("owner_id =?", id).Or("owner_id = 0").Find(&categories)

	return c.JSON(categories)
}

func AddCategoryByUser(c fiber.Ctx) error {
	logging.Logger.Info("Request to add category")

	userId, err2, isCheck := CheckUser(c)
	if isCheck {
		return err2
	}

	var data map[string]string
	if err := c.Bind().Body(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}
	if data["name"] == "" || data["description"] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}
	existingCategory := models.Category{}
	if err := database.DB.Where("name = ?", data["name"]).
		Where("owner_id = ? OR owner_id = 0", userId).First(&existingCategory).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category already exists",
		})
	}

	var category models.Category
	category.Name = data["name"]
	category.Description = data["description"]
	category.OwnerId = userId
	if err := database.DB.Create(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create category",
		})
	}

	return c.JSON(category)
}
