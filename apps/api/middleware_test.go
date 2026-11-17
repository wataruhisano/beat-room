package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWithRequestLogging(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /rooms/{token}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	req := httptest.NewRequest(http.MethodGet, "/rooms/secret-token", nil)
	rec := httptest.NewRecorder()

	withRequestLogging(logger, mux).ServeHTTP(rec, req)

	requestID := rec.Header().Get("X-Request-Id")
	if requestID == "" {
		t.Fatal("X-Request-Id header is empty")
	}
	if strings.Contains(logs.String(), "secret-token") {
		t.Errorf("log contains the URL token: %s", logs.String())
	}

	var entry map[string]any
	if err := json.Unmarshal(logs.Bytes(), &entry); err != nil {
		t.Fatalf("parse log: %v", err)
	}
	if got := entry["route"]; got != "GET /rooms/{token}" {
		t.Errorf("route = %v, want %q", got, "GET /rooms/{token}")
	}
	if got := entry["status"]; got != float64(http.StatusTeapot) {
		t.Errorf("status = %v, want %d", got, http.StatusTeapot)
	}
	if got := entry["request_id"]; got != requestID {
		t.Errorf("request_id = %v, want %q", got, requestID)
	}
}

func TestWithRecovery(t *testing.T) {
	handler := withRecovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if got := rec.Body.String(); got != `{"status":"error"}` {
		t.Errorf("body = %q, want %q", got, `{"status":"error"}`)
	}
}
