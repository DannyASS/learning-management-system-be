package course_route

import (
	"github.com/DannyAss/users/config"
	"github.com/DannyAss/users/internal/database"
	courses_handler "github.com/DannyAss/users/internal/handler/courses"
	courses_repository "github.com/DannyAss/users/internal/repositories/courses"
	courses_usecase "github.com/DannyAss/users/internal/usecase/courses"
	"github.com/gofiber/fiber/v2"
)

func CourseRotes(router fiber.Router, db *database.DBManager, cfg *config.ConfigEnv) {
	courseRepo := courses_repository.NewCourseRepos(db)
	CourseUsecase := courses_usecase.NewCourseUsecase(db, courseRepo)
	CourseHandler := courses_handler.NewCourseHandler(CourseUsecase)

	Course := router.Group("course")
	{
		Course.Post("/create/", CourseHandler.CreateCourse)
		Course.Get("/:id", CourseHandler.GetCourse)
		Course.Get("/", CourseHandler.GetListCourse)
		Course.Put("/update/:id", CourseHandler.UpdateCourse)
	}

}
