package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/samber/do/v2"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/config"
	"github.com/ulbwa/telegram-oidc-provider/internal/infrastructure/di"
	"gorm.io/gorm"
)

const gracefulShutdownTimeout = 15 * time.Second

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Read(*configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	injector := di.NewContainer(cfg)

	logger, err := do.Invoke[zerolog.Logger](injector)
	if err != nil {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	// Replace process-wide default loggers with DI-provided logger.
	zlog.Logger = logger
	stdlog.SetFlags(0)
	stdlog.SetOutput(logger)

	app, err := do.Invoke[*echo.Echo](injector)
	if err != nil {
		return fmt.Errorf("failed to build echo app: %w", err)
	}

	serverErrCh := make(chan error, 1)
	go func() {
		logger.Info().Str("address", cfg.HTTPServer.Address).Msg("starting http server")
		if err := app.Start(cfg.HTTPServer.Address); err != nil {
			serverErrCh <- err
			return
		}
		serverErrCh <- nil
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	select {
	case err := <-serverErrCh:
		if err != nil {
			return fmt.Errorf("http server error: %w", err)
		}
		return nil
	case sig := <-sigCh:
		logger.Info().Str("signal", sig.String()).Msg("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	var shutdownErrs []error

	if err := app.Shutdown(shutdownCtx); err != nil {
		shutdownErrs = append(shutdownErrs, fmt.Errorf("echo shutdown failed: %w", err))
	}

	if db, err := do.Invoke[*gorm.DB](injector); err == nil {
		sqlDB, dbErr := db.DB()
		if dbErr != nil {
			shutdownErrs = append(shutdownErrs, fmt.Errorf("db handle error: %w", dbErr))
		} else if closeErr := sqlDB.Close(); closeErr != nil {
			shutdownErrs = append(shutdownErrs, fmt.Errorf("db close failed: %w", closeErr))
		}
	}

	if redisClient, err := do.Invoke[*redis.Client](injector); err == nil {
		if closeErr := redisClient.Close(); closeErr != nil {
			shutdownErrs = append(shutdownErrs, fmt.Errorf("redis close failed: %w", closeErr))
		}
	}

	if err := errors.Join(shutdownErrs...); err != nil {
		return err
	}

	logger.Info().Msg("graceful shutdown completed")
	return nil
}
