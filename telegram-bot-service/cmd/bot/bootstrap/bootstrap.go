package bootstrap

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
	"telegram-bot-service/cmd/bot/handlers"
	"telegram-bot-service/cmd/bot/middleware"
	"telegram-bot-service/internal/platform/client"
)

type TelegramBotGroup struct {
}

func NewTelegramBotGroup() TelegramBotGroup {
	return TelegramBotGroup{}
}

func Start(
	lc fx.Lifecycle,
	gdBotHandler *handlers.GDFamilyUnityBot,
	usersClient *client.UserServiceClient,
	logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info(fmt.Sprintf("Starting Telegram Bots Hanlers..."))

			go func() {
				logger.Info(fmt.Sprintf("gdBotHandler..."))

				usernames, err := usersClient.GetUsernames()
				if err != nil {
					logger.Error(err.Error())
					return
				}

				gdBotHandler.TelegramBot.Use(middleware.Authorize(logger, usernames))
				gdBotHandler.TelegramBot.Handle(telebot.OnPhoto, gdBotHandler.UploadImage)
				gdBotHandler.TelegramBot.Start()
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing server...")

			return nil
		},
	})
}
