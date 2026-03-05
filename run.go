package main

import (
	"net/http"
)

func run() error {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	fsHandler := http.FileServer(http.Dir("."))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		// Set the content-type header
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.Handle("/app/", http.StripPrefix("/app", fsHandler))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	return srv.ListenAndServe()
}
