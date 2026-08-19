package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"__GO_MODULE__/api/internal/httpx"
	"github.com/jackc/pgx/v5/pgxpool"
)

func routes(pool *pgxpool.Pool, loggers ...*slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", health)
	mux.HandleFunc("GET /api/v1/health", health)
	mux.HandleFunc("GET /readyz", readiness(pool))
	logger := slog.Default()
	if len(loggers) > 0 && loggers[0] != nil {
		logger = loggers[0]
	}
	withIdempotency := httpx.IdempotencyMiddleware(httpx.NewMemoryIdempotencyStore(), mux)
	return httpx.RequestIDMiddleware(httpx.LoggingMiddleware(logger, withIdempotency))
}

func readiness(pool *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if pool == nil {
			writeJSON(writer, http.StatusOK, map[string]string{"status": "ok", "database": "not configured"})
			return
		}
		if err := pool.Ping(request.Context()); err != nil {
			writeJSON(writer, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
			return
		}
		writeJSON(writer, http.StatusOK, map[string]string{"status": "ok", "database": "ok"})
	}
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}

func health(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(map[string]string{"status": "ok"})
}
