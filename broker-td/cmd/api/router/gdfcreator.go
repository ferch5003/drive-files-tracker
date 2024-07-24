package router

import (
	"broker-td/cmd/api/handler"
	"broker-td/config"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var NewGDriveFolderCreatorModule = fx.Module("gdrive-folder-creator",
	// Register Handler
	fx.Provide(handler.NewGDriveFolderCreatorHandler),

	// Register Router
	fx.Provide(
		fx.Annotate(
			NewGDriveFolderCreatorRouter,
			fx.ResultTags(`group:"routers"`),
		),
	),
)

type gDriveFolderCreatorRouter struct {
	App     fiber.Router
	config  *config.EnvVars
	Handler *handler.GDriveFolderCreatorHandler
}

func NewGDriveFolderCreatorRouter(
	app *fiber.App,
	config *config.EnvVars,
	gdfCreatorHandler *handler.GDriveFolderCreatorHandler) Router {
	return &gDriveFolderCreatorRouter{
		App:     app,
		config:  config,
		Handler: gdfCreatorHandler,
	}
}

func (g gDriveFolderCreatorRouter) Register() {
	g.App.Route("/gdrive-folder-creator", func(api fiber.Router) {
		api.Post("/", g.Handler.Post).Name("post")
	}, "gdrive-folder-creator.")
}
