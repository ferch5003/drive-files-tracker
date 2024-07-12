package router

import (
	"broker-td/cmd/api/handler"
	"broker-td/config"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var NewUserModule = fx.Module("gdrive-family-uploader",
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

func NewUserRouter(
	app *fiber.App,
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
	}, "users.")
}
