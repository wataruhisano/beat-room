package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, r, http.StatusOK, statusResponse{Status: "ok"})
}

type pinger interface {
	Ping(ctx context.Context) error
}

func readyz(db pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			loggerFrom(r.Context()).Error("ping database", "error", err)
			writeJSON(w, r, http.StatusServiceUnavailable, statusResponse{Status: "unavailable"})
			return
		}

		writeJSON(w, r, http.StatusOK, statusResponse{Status: "ok"})
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		logger.Error("DATABASE_URL is required")
		os.Exit(1)
	}

	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		logger.Error("create database pool", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz)
	mux.HandleFunc("GET /readyz", readyz(db))

	handler := withRequestLogging(logger, withRecovery(mux))

	logger.Info("listening", "port", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
