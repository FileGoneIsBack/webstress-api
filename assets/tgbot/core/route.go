package handler

import (
	"net/http"
	api1 "bot/core/api" 
)

func NewAPIHandler() http.Handler {
	mux := http.NewServeMux()
	api1.RegisterRoutes(mux)
	api1.RegisterRoute2(mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "API Server is running"}`))
	})

	return mux
}