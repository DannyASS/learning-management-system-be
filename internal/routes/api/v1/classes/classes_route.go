package classes_route

import (
	"github.com/DannyAss/users/config"
	"github.com/DannyAss/users/internal/database"
	classes_handler "github.com/DannyAss/users/internal/handler/classes"
	classes_repository "github.com/DannyAss/users/internal/repositories/classes"
	courses_repository "github.com/DannyAss/users/internal/repositories/courses"
	import_repository "github.com/DannyAss/users/internal/repositories/import"
	classes_usecase "github.com/DannyAss/users/internal/usecase/classes"
	"github.com/gofiber/fiber/v2"
)

func ClassesRoute(router fiber.Router, db *database.DBManager, cfg *config.ConfigEnv) {
	repos := classes_repository.NewClassesRepos(db)
	Importrepos := import_repository.NewImportRepository(db)
	Crepos := courses_repository.NewCourseRepos(db)
	useCase := classes_usecase.NewClassesUsecase(repos, Importrepos, db, Crepos)
	handler := classes_handler.NewClasssesHandler(useCase)

	class := router.Group("/class")
	{
		class.Get("/", handler.GetlistClassPage)
		class.Get("/:id", handler.GetClassByIDClassHdr)
		class.Put("/:id", handler.UpdateClassHdr)
		class.Post("/", handler.CreateClassHdr)
		class.Get("/template/download", handler.DownloadTemplate)
	}

	student := router.Group("/student")
	{
		student.Get("/class/:id", handler.GetStudentClassByIDClass)
		student.Delete("/:id", handler.DeleteStudentClassByID)
		student.Delete("/class/:id", handler.DeleteStudentClassByIDClass)
	}

	courses := router.Group("/courses")
	{
		courses.Get("/template/download", handler.DownloadTemplateCourse)
		courses.Get("/class/:id", handler.GetCOurseClassByIDClass)
	}

	dashboard := router.Group("/dashboard")
	{
		dashboard.Get("/class/inform/:id", handler.GetInformDasboardClass)
		dashboard.Get("/class/module/:id", handler.GetAllModulByClassAndRole)
		dashboard.Get("/class/course/:id", handler.GetAvailablecourse)
		dashboard.Get("/class/module/avail/:id", handler.GetAvailableModulDash)
		dashboard.Post("/class/module/add/:id", handler.AddModulDash)
		dashboard.Put("/class/module/update/:id", handler.UpdateModulDash)
		dashboard.Get("/class/template/absen/download/:id", handler.DownloadTemplateAbsen)
	}
}
