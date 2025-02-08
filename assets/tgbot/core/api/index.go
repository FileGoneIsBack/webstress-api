package api

import (
	"net/http"
	"fmt"
)

func RegisterRoute2(mux *http.ServeMux) {
	mux.HandleFunc("/api2", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"api1": "Hello from API 1"}`)
	})

}