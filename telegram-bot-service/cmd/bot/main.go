package main

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
	"telegram-bot-service/cmd/bot/bootstrap"
	"telegram-bot-service/cmd/bot/handlers"
	"telegram-bot-service/config"
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

	app := fx.New(
		// creates: config.EnvVars
		fx.Supply(configurations),
		// creates: *zap.Logger
		fx.Supply(logger),

		// creates: *bootstrap.TelegramBotGroup
		fx.Provide(bootstrap.NewTelegramBotGroup),

		// creates: *handlers.GDFamilyUnityBot
		fx.Provide(handlers.NewGDFamilyUnityBot),

		// Start Bots.
		fx.Invoke(bootstrap.Start),
	)

	app.Run()
}
