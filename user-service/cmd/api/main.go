package main

import (
	"context"
	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"log"
	"time"
	"user-service/cmd/api/bootstrap"
	appcrons "user-service/cmd/api/cron"
	"user-service/cmd/api/router"
	"user-service/config"
	"user-service/internal/botuser"
	"user-service/internal/platform/client"
	"user-service/internal/platform/postgresql"
)

func main() {
	configurations, err := config.NewConfigurations()
	if err != nil {
		log.Fatalln(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln(err)
	}

	location, err := time.LoadLocation("America/Bogota")
	if err != nil {
		log.Fatalln(err)
	}

	brokerClient := client.NewBrokerClient(configurations.BrokerTDBaseURL)

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
		fx.Provide(context.Background),

		// creates: *sqlx.DB
		fx.Provide(postgresql.NewConnection),

		// creates: *client.BrokerClient
		fx.Supply(brokerClient),

		// Provide modules
		router.NewUserModule,

		// creates: *botuser.Repository
		fx.Provide(botuser.NewRepository),
		// creates: *botuser.Service
		fx.Provide(botuser.NewService),

		// creates: *cron.Cron
		fx.Supply(cron.New(
			cron.WithLocation(location))),

		// creates: *cron.FolderCronJob
		fx.Provide(appcrons.NewFolderCronJob),

		// Start web server.
		fx.Invoke(bootstrap.Start),
	)

	app.Run()
}
