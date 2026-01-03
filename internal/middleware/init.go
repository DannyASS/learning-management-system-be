package middleware

import (
	"github.com/DannyAss/users/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func InitMiddlewares(app *fiber.App, cfg *config.ConfigEnv) {
	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.AppEnv != "prod",
	}))
	app.Use(logger.New())
	InitCORS(app, cfg)
	app.Use(etag.New()) // Cache-friendly for GET requests
	// app.Use(limiter.New()) // Optional: Basic rate-limiting middleware
}
