package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/himbo22/source-base/internal/config"
	"github.com/himbo22/source-base/internal/di"

	"go.uber.org/zap"
)

func main() {
	loadConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app, cleanup, err := di.InitApp(loadConfig)
	if err != nil {
		log.Fatalf("Failed to init app: %v", err)
	}
	defer cleanup()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	s := http.Server{Addr: ":" + loadConfig.Server.Port, Handler: app.EchoApp}
	// Start server
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Logger.Error("failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		app.Logger.Error("failed to stop server", zap.Error(err))
	}
}
