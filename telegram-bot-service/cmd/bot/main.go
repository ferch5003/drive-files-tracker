package main

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"telegram-bot-service/cmd/bot/bootstrap"
	"telegram-bot-service/cmd/bot/handlers"
	"telegram-bot-service/config"
	"telegram-bot-service/internal/platform/client"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	configurations, err := config.NewConfigurations()
	if err != nil {
		logger.Fatal(err.Error())
	}

	usersClient := client.NewUserServiceClient(configurations.BrokerTDBaseURL)

	app := fx.New(
		// creates: config.EnvVars
		fx.Supply(configurations),
		// creates: *zap.Logger
		fx.Supply(logger),

		// creates: *client.UserServiceClient
		fx.Supply(usersClient),

		// creates: *bootstrap.TelegramBotGroup
		fx.Provide(bootstrap.NewTelegramBotGroup),

		// creates: *handlers.GDFamilyUnityBot
		fx.Provide(handlers.NewGDFamilyUnityBot),
		// creates: *handlers.GDFamilyGardenBot
		fx.Provide(handlers.NewGDFamilyGardenBot),
		// creates: *handlers.GDOSCommercialBot
		fx.Provide(handlers.NewGDOSCommercialBot),

		// Start Bots.
		fx.Invoke(bootstrap.Start),
	)

	app.Run()
}
