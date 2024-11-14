package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"project/database"
	"project/logging"
	"project/models"
	"strconv"
	"time"
)

func GetExpenses(c fiber.Ctx) error {
	logging.Logger.Info("Request to get expenses")

	id, err2, done := CheckUser(c)
	if done {
		return err2
	}
	var expenses []models.Expense
	database.DB.Where("user_id =?", id).Find(&expenses)

	return c.JSON(expenses)
}

func AddExpenseByUser(c fiber.Ctx) error {
	logging.Logger.Info("Request to add expense")

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
	if data["name"] == "" || data["category_id"] == "" || data["amount"] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	var category models.Category
	if err := database.DB.Where("id = ?", data["category_id"]).Where("owner_id = ? OR owner_id = 0", userId).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Category not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	var expense models.Expense
	expense.Name = data["name"]
	expense.UserID = userId
	expense.CategoryID = category.ID
	amountStr := data["amount"]
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid amount format",
		})
	}
	expense.Amount = amount
	if dateStr, ok := data["date"]; ok && dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date format",
			})
		}
		expense.Date = parsedDate
	} else {
		expense.Date = time.Now()
	}

	if err := database.DB.Create(&expense).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create expense",
		})
	}

	return c.JSON(expense)
}

func DeleteExpense(c fiber.Ctx) error {
	logging.Logger.Info("Request to delete expense")

	id, err2, done := CheckUser(c)
	if done {
		return err2
	}
	idStr := c.Params("id")
	expenseId, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid expense ID",
		})
	}
	var expense models.Expense
	if err := database.DB.Where("id = ?", expenseId).Where("user_id = ?", id).First(&expense).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Expense not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	if err := database.DB.Delete(&expense).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete expense",
		})
	}
	return c.JSON(fiber.Map{
		"message": "Expense deleted successfully",
	})
}

func UpdateExpense(c fiber.Ctx) error {
	logging.Logger.Info("Request to update expense")
	id, err2, done := CheckUser(c)
	if done {
		return err2
	}
	idStr := c.Params("id")
	expenseId, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid expense ID",
		})
	}
	var expense models.Expense
	if err := database.DB.Where("id = ?", expenseId).Where("user_id = ?", id).First(&expense).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Expense not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	var data map[string]string
	if err := c.Bind().Body(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}
	if data["name"] != "" {
		expense.Name = data["name"]
	}
	if data["category_id"] != "" {
		var category models.Category
		if err := database.DB.Where("id = ?", data["category_id"]).Where("owner_id = ? OR owner_id = 0", id).First(&category).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Category not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
		expense.CategoryID = category.ID
	}
	if data["amount"] != "" {
		amountStr := data["amount"]
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid amount format",
			})
		}
		expense.Amount = amount
	}
	if data["date"] != "" {
		parsedDate, err := time.Parse("2006-01-02", data["date"])
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date format",
			})
		}
		expense.Date = parsedDate
	}
	if err := database.DB.Save(&expense).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update expense",
		})
	}
	return c.JSON(expense)
}

func GetSumExpensesByCategoryId(c fiber.Ctx) error {
	logging.Logger.Info("Request to get sum expenses by category")

	id, err2, done := CheckUser(c)
	if done {
		return err2
	}
	idStr := c.Params("category_id")
	categoryId, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}
	var sum float64
	database.DB.Model(&models.Expense{}).Where("category_id = ?", categoryId).Where("user_id = ?", id).Select("SUM(amount)").Row().Scan(&sum)
	return c.JSON(fiber.Map{
		"sum": sum,
	})
}

func GetSumExpenses(c fiber.Ctx) error {
	logging.Logger.Info("Request to get sum expenses")

	id, err2, done := CheckUser(c)
	if done {
		return err2
	}
	var sum float64
	database.DB.Model(&models.Expense{}).Where("user_id = ?", id).Select("SUM(amount)").Row().Scan(&sum)
	return c.JSON(fiber.Map{
		"sum": sum,
	})
}
