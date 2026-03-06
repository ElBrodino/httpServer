package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverhits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverhits.Add(1)
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

	mux.HandleFunc("GET /api/healthz",
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

	mux.HandleFunc("GET /admin/metrics",
		func(w http.ResponseWriter, r *http.Request) {
			hits := cfg.fileserverhits.Load()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)

			metricsString := fmt.Appendf([]byte(``), `
				<html>
					<body>
						<h1>Welcome, Chirpy Admin</h1>
						<p>Chirpy has been visited %d times!</p>
					<body>
				</html>
				`, hits)

			w.Write(metricsString)
		})

	mux.HandleFunc("POST /admin/reset",
		func(w http.ResponseWriter, r *http.Request) {
			cfg.fileserverhits.Store(0)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Reset hits counter"))
		})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return srv.ListenAndServe()
}
