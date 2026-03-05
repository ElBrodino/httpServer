package main

import (
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
)

type apiConfig struct {
	fileserverhits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverhits.Add(1)
		log.Println("counted:", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func run() error {
	const filepathRoot = "."
	const port = "8080"

	cfg := &apiConfig{}
	mux := http.NewServeMux()

	fsHandler := http.FileServer(http.Dir("."))
	appHandler := http.StripPrefix("/app", fsHandler)
	mux.Handle("/app/", cfg.middlewareMetricsInc(appHandler))

	healthzHandler := func(w http.ResponseWriter, req *http.Request) {
		// Set the content-type header
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}

	metricsHandler := func(w http.ResponseWriter, r *http.Request) {
		hits := cfg.fileserverhits.Load()
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits: " + strconv.Itoa(int(hits))))
	}

	resetHandler := func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverhits.Store(0)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reset hits counter"))
	}

	mux.HandleFunc("GET /healthz", healthzHandler)
	mux.HandleFunc("GET /metrics", metricsHandler)
	mux.HandleFunc("POST /reset", resetHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return srv.ListenAndServe()
}
