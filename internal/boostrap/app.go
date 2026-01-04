package boostrap

import (
	"log"
	"os"
	"runtime"

	"github.com/DannyAss/users/config"
	"github.com/DannyAss/users/internal/database"
	"github.com/DannyAss/users/internal/middleware"
	router "github.com/DannyAss/users/internal/routes/api"
	"github.com/DannyAss/users/internal/worker"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

func Buildapp(cfg *config.ConfigEnv) (*fiber.App, func(), error) {
	prefork := runtime.GOOS != "windows" && cfg.AppEnv != "dev"
	dbmanager := database.NewDBManager(cfg.DBConnnect)

	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered panic in worker PID", os.Getpid(), ":", r)
		}
	}()

	engine := html.New("./internal/resources/templates", ".html")
	engine.AddFunc("T", func(key string) string {
		return key
	})
	engine.Debug(cfg.AppDebug)

	app := fiber.New(fiber.Config{
		Prefork:       prefork,
		StrictRouting: false,
		CaseSensitive: false,
		ServerHeader:  cfg.AppName,
		JSONEncoder:   sonic.Marshal,
		JSONDecoder:   sonic.Unmarshal,
		Views:         engine,
		// BodyLimit: 8 * 1024 * 1024, // 8MB
	})

	app.Use(logger.New())

	middleware.InitMiddlewares(app, cfg)

	worker.InitJobQueue(1000)
	worker.StartWorkers(3)

	router.InitAllRoutes(app, dbmanager, cfg)

	// Cleanup function to close db connections gracefully
	cleanup := func() {
		_ = app.Shutdown()
		dbmanager.Close()
		log.Println("Shutdown completed: app & database closed")
	}

	return app, cleanup, nil
}
