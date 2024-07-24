package main

import (
	"broker-td/cmd/api/bootstrap"
	"broker-td/cmd/api/router"
	"broker-td/config"
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"log"
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
		router.NewGDriveFamilyUploaderModule,
		router.NewGDriveFolderCreatorModule,
		router.NewUserModule,

		// Start web server.
		fx.Invoke(bootstrap.Start),
	)

	app.Run()
}
