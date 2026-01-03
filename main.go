package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DannyAss/users/config"
	"github.com/DannyAss/users/internal/boostrap"
)

func main() {
	cfg := config.InitConfig()
	app, cleanup, err := boostrap.Buildapp(cfg)
	if err != nil {
		log.Fatalf("Error initializing app gofiber: %v", err)
	}
	defer cleanup()

	addr := cfg.AppPort
	if addr != "" && addr[0] != ':' {
		addr = ":" + addr
	}

	go func() {
		log.Printf("listening on %s\n", cfg.AppURL+addr)
		if err := app.Listen(addr); err != nil {
			log.Printf("App gofiber stopped: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}
