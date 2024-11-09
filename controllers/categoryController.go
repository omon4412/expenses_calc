package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"log"
	"project/config"
	"project/database"
	"project/logging"
	"project/models"
	"strconv"
)

func GetCategories(c fiber.Ctx) error {
	logging.Logger.Info("Request to get categories")

	id, err2, done := checkUser(c)
	if done {
		return err2
	}
	var categories []models.Category
	database.DB.Where("owner_id =?", id).Or("owner_id = 0").Find(&categories)

	return c.JSON(categories)
}

func AddCategoryByUser(c fiber.Ctx) error {
	logging.Logger.Info("Request to get categories")

	userId, err2, isCheck := checkUser(c)
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

func checkUser(c fiber.Ctx) (uint, error, bool) {
	cookie := c.Cookies("jwt")

	secretKey := config.GetConfig().JWT.Secret
	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		}), true
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return 0, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse claims",
		}), true
	}

	id, _ := strconv.Atoi((*claims)["sub"].(string))
	user := models.User{ID: uint(id)}

	if err := database.DB.Where("id = ?", user.ID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			}), true
		}
		log.Println("Database error:", err)
		return 0, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving user",
		}), true
	}
	return uint(id), nil, false
}
