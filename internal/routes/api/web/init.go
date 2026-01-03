package web_routes

import (
	"github.com/DannyAss/users/config"
	"github.com/DannyAss/users/internal/database"
	web_home_handler "github.com/DannyAss/users/internal/handler/web/home"

	"github.com/gofiber/fiber/v2"
)

func InitWebRoutes(app *fiber.App, db *database.DBManager, cfg *config.ConfigEnv) {
	// Initialize handlers
	homeHandler := web_home_handler.NewHomeHandler(cfg)

	// Static files
	app.Static("/public", "./public")

	// Web routes
	app.Get("/", homeHandler.HomePage)
}
