package router

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"user-service/cmd/api/handler"
	"user-service/config"
	"user-service/internal/bot"
	"user-service/internal/user"
)

var NewUserModule = fx.Module("user",
	// Register Repository & Service
	fx.Provide(user.NewRepository),
	fx.Provide(bot.NewRepository),
	fx.Provide(user.NewService),

	// Register Handler
	fx.Provide(handler.NewUserHandler),

	// Register Router
	fx.Provide(
		fx.Annotate(
			NewUserRouter,
			fx.ResultTags(`group:"routers"`),
		),
	),
)

type userRouter struct {
	App     fiber.Router
	config  *config.EnvVars
	Handler *handler.UserHandler
}

func NewUserRouter(app *fiber.App,
	config *config.EnvVars,
	userHandler *handler.UserHandler) Router {
	return &userRouter{
		App:     app,
		config:  config,
		Handler: userHandler,
	}
}

func (u userRouter) Register() {
	u.App.Route("/users", func(api fiber.Router) {
		api.Get("/", u.Handler.GetAll).Name("get_all")
		api.Get("/:user_username/bot/:bot_name", u.Handler.FindFolderID).Name("fid_folder_id")
		api.
			Get("/:user_username/bot/:bot_name/date/:date/spreadsheets", u.Handler.GetSpreadsheetData).
			Name("get_spreadsheet_data")
	}, "users.")
}
