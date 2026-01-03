package middleware

import (
	"fmt"
	"log"
	"strings"

	"github.com/DannyAss/users/config"
	"github.com/DannyAss/users/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InitCORS(app *fiber.App, cfg *config.ConfigEnv) {
	origins := strings.Join(utils.CleanText(cfg.CORSAllowedOrigins), ",")
	headers := strings.Join(utils.CleanText(cfg.CORSAllowedHeaders), ",")
	methods := strings.Join(utils.CleanText(cfg.CORSAllowedMethods), ",")

	if cfg.CORSAllowCredentials && origins == "*" {
		log.Println("CORS_ALLOWED_CREDENTIALS=true but origins='*'. Browsers will block. Consider listing explicit origins.")
	}

	fmt.Println("allow credential :", cfg.CORSAllowCredentials)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowHeaders:     headers,
		AllowMethods:     methods,
		AllowCredentials: cfg.CORSAllowCredentials,
	}))
}
