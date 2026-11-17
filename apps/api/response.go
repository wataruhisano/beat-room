package main

import (
	"encoding/json"
	"net/http"
)

type statusResponse struct {
	Status string `json:"status"`
}

func writeJSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	body, err := json.Marshal(v)
	if err != nil {
		loggerFrom(r.Context()).Error("encode response", "error", err)
		http.Error(w, `{"status":"error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(body); err != nil {
		loggerFrom(r.Context()).Error("write response", "error", err)
	}
}
