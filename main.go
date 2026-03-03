package main

import (
	"net/http"
)

func main() {
	const filepathRoot = "."
	sMux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: sMux,
	}
	sMux.Handle("/", http.FileServer(http.Dir(filepathRoot)))
	server.ListenAndServe()
}
