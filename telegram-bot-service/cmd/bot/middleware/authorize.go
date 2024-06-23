package middleware

import (
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
	"slices"
)

var _authorizedUsers = []string{""}

// Authorize only let the authorized users to interact with the application.
func Authorize(logger *zap.Logger) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			if !slices.Contains(_authorizedUsers, c.Sender().Username) {
				logger.Warn(fmt.Sprintf("This user is not authorized: @%s", c.Sender().Username))
				return c.Send("No estas autorizado para mandar mensajes a este bot...")
			}
			return next(c)
		}
	}
}
