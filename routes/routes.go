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
}
