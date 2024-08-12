package bootstrap

import (
	"broker-td/cmd/api/router"
	"broker-td/config"
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"testing"
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
		fx.Provide(NewFiberServer),

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
