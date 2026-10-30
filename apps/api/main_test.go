package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	// 準備
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	// 実行
	healthz(rec, req)

	// 確認
	if rec.Code != 200 {
		t.Errorf("status = %d, want %d", rec.Code, 200)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want %q", got, "application/json")
	}
	if got := rec.Body.String(); got != `{"status":"ok"}` {
		t.Errorf("body = %q, want %q", got, `{"status":"ok"}`)
	}
}
