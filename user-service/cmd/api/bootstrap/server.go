package bootstrap

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"user-service/cmd/api/cron"
	"user-service/cmd/api/router"
	"user-service/config"
)

const (
	_defaultPort = "3000"
)

func NewFiberServer() *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	return app
}

func Start(
	lc fx.Lifecycle,
	cfg *config.EnvVars,
	app *fiber.App,
	router *router.GeneralRouter,
	folderCronJob *cron.FolderCronJob,
	logger *zap.Logger) {
	port := _defaultPort // Default Port
	if cfg != nil && cfg.Port != "" {
		port = cfg.Port
	}

	// Log all requests.
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	// Allow only internal services.
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://broker-td",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(recover.New())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info(fmt.Sprintf("Starting fiber server on 0.0.0.0:%s", port))

			router.Register()

			go func() {
				logger.Info("Starting...")

				if err := app.Listen(":" + port); err != nil {
					logger.Error(err.Error())
				}
			}()

			go func() {
				if !cfg.ActivateCRON {
					logger.Info("CRON jobs disabled, skipping...")

					return
				}

				logger.Info("Starting CRON jobs...")

				if err := folderCronJob.Run(); err != nil {
					logger.Info("err: ", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing server...")

			return app.Shutdown()
		},
	})
}
