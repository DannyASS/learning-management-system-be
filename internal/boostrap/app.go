package boostrap

import (
	"log"
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

	engine := html.New("./internal/resources/templates", ".html")
	engine.AddFunc("T", func(key string) string { return key })
	engine.Debug(cfg.AppDebug)

	app := fiber.New(fiber.Config{
		Prefork:       false,
		StrictRouting: false,
		CaseSensitive: false,
		ServerHeader:  cfg.AppName,
		JSONEncoder:   sonic.Marshal,
		JSONDecoder:   sonic.Unmarshal,
		Views:         engine,
	})

	// Logger
	app.Use(logger.New())

	// Custom middlewares
	middleware.InitMiddlewares(app, cfg)

	// JobQueue global
	worker.InitJobQueue(1000)

	var dbmanager *database.DBManager

	dbmanager = database.NewDBManager(cfg.DBConnnect)
	if dbmanager == nil {
		log.Fatal("DB Manager initialization failed")
	}
	worker.StartWorkers(3)
	// if !prefork {
	// 	// Non-prefork: DB global
	// 	if dbmanager == nil {
	// 		log.Fatal("DB Manager initialization failed")
	// 	}

	// 	// Start workers normal
	// } else {
	// 	// Prefork: DB per worker di StartWorkersPreforkSafe
	// 	worker.StartWorkersPreforkSafe(3, cfg)
	// }

	// Routes
	router.InitAllRoutes(app, dbmanager, cfg)

	// Cleanup function
	cleanup := func() {
		_ = app.Shutdown()
		if dbmanager != nil {
			dbmanager.Close()
		}
		log.Println("Shutdown completed: app & database closed")
	}

	return app, cleanup, nil
}
