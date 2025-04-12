package server

import (
	"bot/core"
	"log"
	"net/http"
)

func Start() {
	http.Handle("/", handler.NewAPIHandler())
	log.Println("Starting server on https://localhost:8080...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start the non-secure server: %v", err)
		}
}
