package server

import (
	"bot/core"
	"bot/models"
	"log"
	"net/http"
)

func Start() {
	http.Handle("/", handler.NewAPIHandler())

	certFile := "path/to/cert.pem"
	keyFile := "path/to/key.pem"

	log.Println("Starting server on https://localhost:8080...")
	if models.Config.Secure {
		if err := http.ListenAndServeTLS(":443", certFile, keyFile, nil); err != nil {
			log.Fatalf("Failed to start the secure server: %v", err)
		}
	} else {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Failed to start the non-secure server: %v", err)
		}
	}
}
