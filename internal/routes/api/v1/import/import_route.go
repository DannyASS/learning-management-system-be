package import_route

import (
	"github.com/DannyAss/users/config"
	"github.com/DannyAss/users/internal/database"
	import_hanlder "github.com/DannyAss/users/internal/handler/import"
	import_repository "github.com/DannyAss/users/internal/repositories/import"
	import_usecase "github.com/DannyAss/users/internal/usecase/import"
	"github.com/gofiber/fiber/v2"
)

func ImportRoute(router fiber.Router, db *database.DBManager, cfg *config.ConfigEnv) {
	repo := import_repository.NewImportRepository(db)
	usecase := import_usecase.NewImportUsecase(repo)
	handler := import_hanlder.NewImportHandler(usecase)

	Import := router.Group("/import")
	{
		Import.Get("/", handler.GetImportData)
		Import.Get("/status", handler.GetTatusImportByClassId)
	}
}
