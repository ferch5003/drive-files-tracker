package bootstrap

import (
	"broker-td/cmd/api/router"
	"broker-td/config"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
		AllowOrigins: "http://telegram-bot-service, https://user-service",
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

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing server...")

			return app.Shutdown()
		},
	})
}
