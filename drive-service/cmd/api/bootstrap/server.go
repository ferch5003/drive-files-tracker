package bootstrap

import (
	"context"
	"drive-service/cmd/api/rpc"
	"drive-service/config"
	"drive-service/internal/platform/driveaccount"
	"fmt"
	"github.com/otiai10/gosseract/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	ntrpc "net/rpc"
)

const (
	_defaultRPCPort = "5001"
)

func NewServer(ctx context.Context, serviceAccount driveaccount.ServiceAccount) (*rpc.Server, error) {
	client, err := serviceAccount.Get()
	if err != nil {
		return nil, err
	}

	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	sheetService, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return &rpc.Server{
		ServiceAccount: serviceAccount,
		DriveService:   driveService,
		SheetService:   sheetService,
		OCRClient:      gosseract.NewClient(),
	}, nil
}

func Start(
	lc fx.Lifecycle,
	cfg *config.EnvVars,
	rpcServer *rpc.Server,
	logger *zap.Logger) {
	port := _defaultRPCPort // Default Port
	if cfg != nil && cfg.Port != "" {
		port = cfg.Port
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info(fmt.Sprintf("Starting RPC Server on 0.0.0.0:%s", port))

			go func() {
				defer func(OCRClient *gosseract.Client) {
					err := OCRClient.Close()
					if err != nil {
						logger.Error(err.Error())
						return
					}
				}(rpcServer.OCRClient)
				
				// Register the RPC Server
				if err := ntrpc.Register(rpcServer); err != nil {
					logger.Error(err.Error())
					return
				}

				err := rpcServer.Listen(port)
				if err != nil {
					logger.Error(err.Error())
					return
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing server...")

			return nil
		},
	})
}
