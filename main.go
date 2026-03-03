package main

import (
	"net/http"
)

func main() {
	rootString := http.Dir("/")
	sMux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: sMux,
	}
	fileHandler := http.FileServer(rootString)
	sMux.Handle(string(rootString), fileHandler)
	server.ListenAndServe()
}
