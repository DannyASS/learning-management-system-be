package api_v1

import (
	"github.com/DannyAss/users/config"
	"github.com/DannyAss/users/internal/database"
	"github.com/DannyAss/users/internal/middleware"
	classes_route "github.com/DannyAss/users/internal/routes/api/v1/classes"
	course_route "github.com/DannyAss/users/internal/routes/api/v1/course"
	import_route "github.com/DannyAss/users/internal/routes/api/v1/import"
	module_route "github.com/DannyAss/users/internal/routes/api/v1/modules"
	users_route "github.com/DannyAss/users/internal/routes/api/v1/users"
	"github.com/gofiber/fiber/v2"
)

func InitRouteAPI(db *database.DBManager, app *fiber.App, cfg *config.ConfigEnv) {
	api_v1_routes := app.Group("/api/v1")

	users_route.UserRotes(api_v1_routes, db, cfg)
	NeedAuth(db, api_v1_routes, cfg)
}

func NeedAuth(db *database.DBManager, app fiber.Router, cfg *config.ConfigEnv) {
	api_v1_routes := app.Group("", middleware.AuthMiddleware(cfg))
	course_route.CourseRotes(api_v1_routes, db, cfg)
	module_route.ModuleRoute(api_v1_routes, db, cfg)
	classes_route.ClassesRoute(api_v1_routes, db, cfg)
	import_route.ImportRoute(api_v1_routes, db, cfg)
}
