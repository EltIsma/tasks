package app

import (
	"log/slog"
	"os"
	"task/internal/adapters/brokers/kafka"
	"task/internal/adapters/pgrepo"
	"task/internal/adapters/redis"
	"task/internal/config"
	httpserver "task/internal/ports/httpServer"
	"task/internal/services"
	"task/pkg/database"
)

type App struct {
	Server   *httpserver.Server
	Postgres *database.Postgres
	Redis    *redis.Redis
}

func InitApp(cfg *config.Config, logger *slog.Logger) (*App, error) {
	postgres, err := database.NewPG(cfg.Postgres.PostgresURL)
	if err != nil {
		return nil, err
	}

	rds, limiter, err := redis.New(cfg.Redis.Hosts, cfg.Redis.Password, logger)
	if err != nil {
		return nil, err
	}

	kafkaProducer, err := kafka.NewProducer(&cfg.Kafka, logger)
	if err != nil {
		return nil, err
	}

	taskService := services.New(logger, pgrepo.NewRepositoruPG(postgres.GetConn()), rds, kafkaProducer)

	httpServer, err := httpserver.NewHTTPServer(&cfg.Server, logger, taskService, limiter)
	if err != nil {
		return nil, err
	}

	return &App{
		Server:   httpServer,
		Postgres: postgres,
		Redis:    rds,
	}, nil

}

func (a *App) Shutdown() {
	a.Server.Stop()
	a.Postgres.Close()
	a.Redis.Close()
}

func InitLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)

	logger := slog.New(handler)
	return logger
}
