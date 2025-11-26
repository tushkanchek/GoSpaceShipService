package main

import (
	"context"
	"fmt"
	"os/signal"
	"platform/pkg/closer"
	"platform/pkg/logger"
	"syscall"
	"time"


	"payment/internal/app"
	"payment/internal/config"

	"go.uber.org/zap"
)

const (
	configPath = "./deploy/compose/payment/.env"
)

func main() {
	// Load config
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed load config: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "failed to create app: %v", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "failed to run gRPC server: %v", zap.Error(err))
		return
	}
}


func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err:=closer.CloseAll(ctx);err!=nil {
		logger.Error(ctx, "❌ error during Shutdown process", zap.Error(err))
	}
}