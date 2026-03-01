package main

import "net/http"

func main() {
	sMux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: sMux,
	}
	server.ListenAndServe()
}
