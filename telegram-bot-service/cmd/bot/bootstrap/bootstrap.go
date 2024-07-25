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
	gdFamilyUnityBotHandler *handlers.GDFamilyUnityBot,
	gdFamilyGardenBotHandler *handlers.GDFamilyGardenBot,
	gdOSBotHandler *handlers.GDOSCommercialBot,
	usersClient *client.UserServiceClient,
	logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info(fmt.Sprintf("Starting Telegram Bots Hanlers..."))

			usernames, err := usersClient.GetUsernames()
			if err != nil {
				logger.Error(err.Error())
				return err
			}

			go func() {
				logger.Info(fmt.Sprintf("gdFamilyUnityBotHandler..."))

				gdFamilyUnityBotHandler.TelegramBot.Use(middleware.Authorize(logger, usernames))
				gdFamilyUnityBotHandler.TelegramBot.Handle(telebot.OnPhoto, gdFamilyUnityBotHandler.UploadImage)
				gdFamilyUnityBotHandler.TelegramBot.Start()
			}()

			go func() {
				logger.Info(fmt.Sprintf("gdFamilyGardenBotHandler..."))

				gdFamilyGardenBotHandler.TelegramBot.Use(middleware.Authorize(logger, usernames))
				gdFamilyGardenBotHandler.TelegramBot.Handle(telebot.OnPhoto, gdFamilyGardenBotHandler.UploadImage)
				gdFamilyGardenBotHandler.TelegramBot.Start()
			}()

			go func() {
				logger.Info(fmt.Sprintf("gdOSCommercialBotHandler..."))

				gdOSBotHandler.TelegramBot.Use(middleware.Authorize(logger, usernames))
				gdOSBotHandler.TelegramBot.Handle(telebot.OnPhoto, gdOSBotHandler.UploadImage)
				gdOSBotHandler.TelegramBot.Start()
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing server...")

			return nil
		},
	})
}
