package routes

import (
	"github.com/gofiber/fiber/v3"
	"project/controllers"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", controllers.Hello)
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)
	app.Get("/api/categories", controllers.GetCategories)
	app.Post("/api/categories", controllers.AddCategoryByUser)
	app.Get("/api/expenses", controllers.GetExpenses)
	app.Post("/api/expenses", controllers.AddExpenseByUser)
	app.Delete("/api/expenses/:id", controllers.DeleteExpense)
	app.Put("/api/expenses/:id", controllers.UpdateExpense)
	app.Get("/api/expenses/category/:category_id", controllers.GetSumExpensesByCategoryId)
	app.Get("/api/expenses/sum", controllers.GetSumExpenses)
}
