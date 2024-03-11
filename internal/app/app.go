package app

import (
	"time"

	grpcapp "github.com/LDmitryLD/sso/sso/internal/app/grpc_app"
	"github.com/LDmitryLD/sso/sso/internal/services/auth"
	"github.com/LDmitryLD/sso/sso/internal/storage/sqlite"
	"go.uber.org/zap"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *zap.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
