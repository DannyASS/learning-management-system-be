package users_route

import (
	"github.com/DannyAss/users/config"
	"github.com/DannyAss/users/internal/database"
	user_handler "github.com/DannyAss/users/internal/handler/user"
	users_repository "github.com/DannyAss/users/internal/repositories/users"
	user_usecase "github.com/DannyAss/users/internal/usecase/user"
	"github.com/gofiber/fiber/v2"
)

func UserRotes(router fiber.Router, db *database.DBManager, cfg *config.ConfigEnv) {
	userRepo := users_repository.NewReposUser(db)
	userUsecase := user_usecase.NewUserCaseUser(userRepo, cfg, db)
	userHandler := user_handler.NewUserhandler(userUsecase)

	user := router.Group("user")
	{
		user.Post("/register", userHandler.Register)
		user.Post("/login", userHandler.Login)
		user.Post("/logout", userHandler.Logout)
		user.Get("/refreshtoken", userHandler.RefreshToken)
		user.Get("/getuserole", userHandler.GetUsers)
		user.Get("/getuserneedapprove", userHandler.GetUsersNeedApproval)
		user.Get("/all-teacher", userHandler.GetAllTeacher)
	}

}
