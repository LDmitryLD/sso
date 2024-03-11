package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/LDmitryLD/sso/sso/internal/app"
	"github.com/LDmitryLD/sso/sso/internal/config"
	"github.com/LDmitryLD/sso/sso/internal/infrastructure/logs"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("error load env" + err.Error())
	}

	cfg := config.MustLoad()

	log := logs.NewLogger(*cfg, os.Stderr)

	log.Info("starting application")

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping application", zap.Stringer("oc_signal", sign))

	application.GRPCSrv.Stop()

	log.Info("application stopped")
}
