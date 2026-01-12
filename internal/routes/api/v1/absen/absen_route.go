package absen_route

import (
	"github.com/DannyAss/users/config"
	"github.com/DannyAss/users/internal/database"
	absen_handler "github.com/DannyAss/users/internal/handler/absen"
	absen_repository "github.com/DannyAss/users/internal/repositories/absen"
	absen_usecase "github.com/DannyAss/users/internal/usecase/absen"
	"github.com/gofiber/fiber/v2"
)

func AbsenRoute(router fiber.Router, db *database.DBManager, cfg *config.ConfigEnv) {
	repo := absen_repository.NewAbsenRepos(db)
	uc := absen_usecase.NewAbsenUsecase(repo)
	handler := absen_handler.NewHandler(uc)

	absen := router.Group("/absen")
	{
		absen.Post("/:id", handler.AbsenToday)
		absen.Get("/:id", handler.GetAbsenToday)
	}
}
