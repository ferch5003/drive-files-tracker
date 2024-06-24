package main

import (
	"context"
	"drive-service/cmd/api/bootstrap"
	"drive-service/config"
	"drive-service/internal/platform/driveaccount"
	"drive-service/internal/platform/files"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"log"
)

func main() {
	configurations, err := config.NewConfigurations()
	if err != nil {
		log.Fatalln(err)
	}

	credentials, err := files.GetFile("config/client_secret.json")
	if err != nil {
		log.Fatalln(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln(err)
	}

	app := fx.New(
		// creates: config.EnvVars
		fx.Supply(configurations),
		// creates: *zap.Logger
		fx.Supply(logger),
		// creates: context.Context
		fx.Provide(context.Background),

		// creates: Drive configs
		fx.Supply(credentials),

		// creates: *driveaccount.ServiceAccount
		fx.Provide(driveaccount.NewServiceAccount),

		// creates: *rpc.Server
		fx.Provide(bootstrap.NewServer),

		// Start web server.
		fx.Invoke(bootstrap.Start),
	)

	app.Run()
}
