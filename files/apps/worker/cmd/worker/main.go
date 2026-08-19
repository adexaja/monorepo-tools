package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adexaja/shoebox"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	dsn := env("DATABASE_URL", "postgres://app:app@localhost:5432/app?sslmode=disable")
	queue, err := shoebox.New(shoebox.Options{Storage: shoebox.Postgres, DSN: dsn})
	if err != nil {
		logger.Error("worker queue failed", "error", err)
		os.Exit(1)
	}
	queue.Handle("default", func(_ context.Context, message shoebox.Message) error {
		logger.Info("job received", "message_id", message.ID, "bytes", len(message.Payload))
		return nil
	}, shoebox.HandlerOptions{MaxRetries: 5, Timeout: 30 * time.Second})
	logger.Info("worker started", "queue", "default")
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := queue.Shutdown(shutdownCtx); err != nil {
		logger.Error("worker shutdown failed", "error", err)
	}
	logger.Info("worker stopped")
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
