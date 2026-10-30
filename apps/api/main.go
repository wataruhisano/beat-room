package main

import (
	"log"
	"net/http"
	"os"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz)

	log.Printf("listening on :%s", port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
