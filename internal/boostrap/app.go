package boostrap

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

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

// Buildapp membuat Fiber app, workers, dan DB manager sesuai env
func Buildapp(cfg *config.ConfigEnv) (*fiber.App, func(), error) {
	prefork := runtime.GOOS != "windows" && cfg.AppEnv != "dev"

	engine := html.New("./internal/resources/templates", ".html")
	engine.AddFunc("T", func(key string) string { return key })
	engine.Debug(cfg.AppDebug)

	app := fiber.New(fiber.Config{
		Prefork:       prefork,
		StrictRouting: false,
		CaseSensitive: false,
		ServerHeader:  cfg.AppName,
		JSONEncoder:   sonic.Marshal,
		JSONDecoder:   sonic.Unmarshal,
		Views:         engine,
	})

	// Logger
	app.Use(logger.New())

	// Middlewares
	middleware.InitMiddlewares(app, cfg)

	var dbmanager *database.DBManager

	if !prefork {
		// Non-prefork: DB global & JobQueue global
		dbmanager = database.NewDBManager(cfg.DBConnnect)
		if dbmanager == nil {
			log.Fatal("DB Manager initialization failed")
		}

		worker.InitJobQueue(1000)
		worker.StartWorkers(3)
	} else {
		// Prefork-safe: DB per worker & JobQueue per worker
		worker.StartWorkersPreforkSafe(3, cfg)
	}

	// Routes
	router.InitAllRoutes(app, dbmanager, cfg)

	// Cleanup function
	cleanup := func() {
		log.Println("Shutdown initiated...")

		// Stop workers
		worker.StopWorkers(5 * time.Second)

		// Close global DB (non-prefork)
		if dbmanager != nil {
			dbmanager.Close()
		}

		_ = app.Shutdown()
		log.Println("Shutdown completed: app & database closed")
	}

	// Signal handling untuk Docker graceful shutdown
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		log.Println("Received stop signal")
		cleanup()
		os.Exit(0)
	}()

	return app, cleanup, nil
}
