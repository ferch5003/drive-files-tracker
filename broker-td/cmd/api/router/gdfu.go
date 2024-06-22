package router

import (
	"broker-td/cmd/api/handler"
	"broker-td/config"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

var NewGDriveFamilyUploaderModule = fx.Module("gdrive-family-uploader",
	// Register Handler
	fx.Provide(handler.NewGDriveFamilyHandler),

	// Register Router
	fx.Provide(
		fx.Annotate(
			NewGDriveFamilyUploaderRouter,
			fx.ResultTags(`group:"routers"`),
		),
	),
)

type gDriveFamilyUploaderRouter struct {
	App     fiber.Router
	config  *config.EnvVars
	Handler *handler.GDriveFamilyHandler
}

func NewGDriveFamilyUploaderRouter(app *fiber.App,
	config *config.EnvVars,
	gdfuHandler *handler.GDriveFamilyHandler) Router {
	return &gDriveFamilyUploaderRouter{
		App:     app,
		config:  config,
		Handler: gdfuHandler,
	}
}

func (g gDriveFamilyUploaderRouter) Register() {
	g.App.Route("/gdrive-family-uploader", func(api fiber.Router) {
		api.Post("/", g.Handler.Post).Name("post")
	}, "gdrive-family-uploader.")
}
