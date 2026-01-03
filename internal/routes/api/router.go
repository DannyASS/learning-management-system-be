package router

import (
	"github.com/DannyAss/users/config"
	"github.com/DannyAss/users/internal/database"
	api_v1 "github.com/DannyAss/users/internal/routes/api/v1"
	web_routes "github.com/DannyAss/users/internal/routes/api/web"
	"github.com/DannyAss/users/pkg/i18n"
	"github.com/DannyAss/users/pkg/presentation"
	"github.com/gofiber/fiber/v2"
)

func InitAllRoutes(app *fiber.App, db *database.DBManager, cfg *config.ConfigEnv) {

	api_v1.InitRouteAPI(db, app, cfg)
	web_routes.InitWebRoutes(app, db, cfg)

	// 404 Route Not Found
	app.Use(func(c *fiber.Ctx) error {
		return presentation.Response[any]().
			SetStatus(false).
			SetStatusCode(fiber.StatusNotFound).
			SetStatusText("status not found").
			SetMessage(i18n.T(c, "global.message.7")).
			Json(c)
	})
}
