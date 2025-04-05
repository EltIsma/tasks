package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"task/internal/app"
	"task/internal/config"

	"golang.org/x/sync/errgroup"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	logger := app.InitLogger()

	application, err := app.InitApp(cfg, logger)
	if err != nil {
		logger.Error("bad configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer application.Shutdown()

	eg, ctx := errgroup.WithContext(context.Background())
	sigQuit := make(chan os.Signal, 1)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		return application.Server.Run(ctx)
	})

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case s := <-sigQuit:
			logger.Info("Captured signal", slog.String("signal", s.String()))
			return fmt.Errorf("captured signal: %v", s)
		}
	})

	err = eg.Wait()
	logger.Info("Gracefully shutting down the servers", slog.String("error", err.Error()))
}
