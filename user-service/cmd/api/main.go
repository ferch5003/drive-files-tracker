package main

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"log"
	"user-service/cmd/api/bootstrap"
	"user-service/cmd/api/router"
	"user-service/config"
)

func main() {
	configurations, err := config.NewConfigurations()
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln(err)
	}

	app := fx.New(
		// creates: config.EnvVars
		fx.Supply(configurations),
		// creates: *zap.Logger
		fx.Supply(logger),
		// creates: *fiber.Router
		fx.Provide(
			fx.Annotate(
				router.NewRouter,
				fx.ParamTags( // Equivalent to *fiber.App, config.Envars, []Router `group:"routers"` in constructor
					``,
					``,
					`group:"routers"`),
			),
		),
		// creates: *fiber.App
		fx.Provide(bootstrap.NewFiberServer),
		// creates: context.Context
		fx.Supply(ctx),

		// Provide modules
		router.NewUserModule,

		// Start web server.
		fx.Invoke(bootstrap.Start),
	)

	app.Run()
}
