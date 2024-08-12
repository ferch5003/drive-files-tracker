package bootstrap

import (
	"context"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"testing"
	appcrons "user-service/cmd/api/cron"
	"user-service/cmd/api/router"
	"user-service/config"
	"user-service/internal/botuser"
	"user-service/internal/platform/client"
	"user-service/internal/platform/postgresql"
)

type mockUserRouter struct {
	mock.Mock
}

func (m *mockUserRouter) Register() {
	m.Called()
}

func TestStart_Successful(t *testing.T) {
	// Given
	mur := new(mockUserRouter)
	mur.On("Register")

	app := fx.New(
		fx.Supply(
			fx.Annotate(
				mur,
				fx.As(new(router.Router))),
		),
		fx.Provide(router.NewRouter),
		fx.Provide(zap.NewDevelopment),
		fx.Provide(config.NewConfigurations),
		fx.Provide(client.NewBrokerClient),
		fx.Provide(botuser.NewRepository),
		fx.Provide(botuser.NewService),
		fx.Provide(cron.New),
		fx.Provide(appcrons.NewFolderCronJob),
		fx.Provide(NewFiberServer),
		fx.Provide(context.Background),
		fx.Provide(postgresql.NewConnection),

		fx.Invoke(Start),
	)

	ctx := context.Background()

	// When
	err := app.Start(ctx)
	require.NoError(t, err)

	// Then
	err = app.Stop(ctx)
	require.NoError(t, err)
}
