package module_route

import (
	"github.com/DannyAss/users/config"
	"github.com/DannyAss/users/internal/database"

	modules_handler "github.com/DannyAss/users/internal/handler/modules"
	course_repository "github.com/DannyAss/users/internal/repositories/courses"
	modules_repository "github.com/DannyAss/users/internal/repositories/modules"
	modules_usecase "github.com/DannyAss/users/internal/usecase/modules"

	"github.com/gofiber/fiber/v2"
)

func ModuleRoute(route fiber.Router, db *database.DBManager, cfg *config.ConfigEnv) {
	repo := course_repository.NewCourseRepos(db)
	repo1 := modules_repository.NewModulesRepos(db)

	// Init usecase
	uc := modules_usecase.NewModuleUsecase(db, repo1, repo)

	// Init handler
	handler := modules_handler.NewModuleHandler(uc)

	// Group routes
	course := route.Group("modules")

	course.Post("/create", handler.CreateModule)
	course.Put("/update/:id", handler.UpdateModule)
	course.Get("/", handler.GetListModule)
	course.Get("/course", handler.GetCourse)
	course.Get("/:id", handler.GetModule)

}
