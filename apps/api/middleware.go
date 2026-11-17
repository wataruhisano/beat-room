package main

import (
	"context"
	"crypto/rand"
	"log/slog"
	"net/http"
	"time"
)

type loggerKey struct{}

func loggerFrom(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(status int) {
	s.status = status
	s.ResponseWriter.WriteHeader(status)
}

func withRequestLogging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := rand.Text()
		reqLogger := logger.With("request_id", requestID)

		w.Header().Set("X-Request-Id", requestID)
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		r = r.WithContext(context.WithValue(r.Context(), loggerKey{}, reqLogger))

		next.ServeHTTP(rec, r)

		// URL pathはsecret tokenを含みうるので、route patternだけを記録する
		reqLogger.Info("request",
			"method", r.Method,
			"route", r.Pattern,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			v := recover()
			if v == nil {
				return
			}
			if v == http.ErrAbortHandler {
				panic(v)
			}
			loggerFrom(r.Context()).Error("panic recovered", "panic", v)
			writeJSON(w, r, http.StatusInternalServerError, statusResponse{Status: "error"})
		}()

		next.ServeHTTP(w, r)
	})
}
