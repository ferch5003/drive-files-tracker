package bootstrap

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
	"telegram-bot-service/cmd/bot/handlers"
)

type TelegramBotGroup struct {
}

func NewTelegramBotGroup() TelegramBotGroup {
	return TelegramBotGroup{}
}

func Start(
	lc fx.Lifecycle,
	gdBotHandler *handlers.GDFamilyUnityBot,
	logger *zap.Logger) {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info(fmt.Sprintf("Starting Telegram Bots Hanlers..."))

			go func() {
				logger.Info(fmt.Sprintf("gdBotHandler..."))
				gdBotHandler.TelegramBot.Handle(telebot.OnPhoto, gdBotHandler.GetPhoto)
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
