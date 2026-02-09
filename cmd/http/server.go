package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/do/v2"
	"github.com/ulbwa/telegram-oidc-provider/internal/common"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/di"
)

func main() {
	// Load configuration
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := common.ReadConfig(configPath)
	if err != nil {
		panic(err)
	}

	// Initialize DI container
	injector := di.NewContainer(cfg)

	// Get HTTP Server
	httpServer := do.MustInvoke[*fiber.App](injector)

	// Run Server in goroutine
	go func() {
		if err := httpServer.Listen(cfg.HTTPServer.Address); err != nil {
			panic(err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := injector.Shutdown(); err != nil {
		panic(err)
	}
}
