package main

import (
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(500)
	log.Printf("%v: %d", msg, code)
}
