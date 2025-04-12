package server

import (
	"api/core/models"
	_ "api/core/models/antiflood"
	"api/core/models/log"
	"errors"
	"fmt"
	"net/http"

	"strings"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/http2"
)

type Server struct {
	server *http.Server
	router *mux.Router
	routes map[string]*Route
	config *Config
}

func NewServer(config *Config) *Server {
	s := &Server{
		server: &http.Server{
			Addr:    config.Addr,
			Handler: nil,
			WriteTimeout: 30 * time.Second,
			ReadTimeout:  30 * time.Second,
		},
		router: mux.NewRouter(),
		routes: make(map[string]*Route),
		config: config,
	}
 

	
	// Configure HTTP/2
	http2.ConfigureServer(s.server, &http2.Server{
		MaxConcurrentStreams: 50, 
		IdleTimeout:          30 * time.Minute,
	})

	return s
}

func (s *Server) ListenAndServe() error {

	s.server.Handler = s.router

	if models.Config.Secure {
		cert := models.Config.Cert
		key := models.Config.Key
		if cert == "" || key == "" {
			return errors.New("certificate or key is empty")
		}
		s.server.Addr = strings.Split(s.config.Addr, ":")[0] + ":80"
		log.Println("Server is running on HTTPS on " + s.server.Addr)
		return s.server.ListenAndServeTLS(cert, key)
	} else {
		log.Println("Server is running on HTTP on " + s.server.Addr)
	}

	log.Println("listening with " + fmt.Sprint(s.Subrouters()) + " subrouters and " + fmt.Sprint(s.Routes()) + " routes.")

	return s.server.ListenAndServe()
}

func (s *Server) Subrouters() int {
	subs := 0
	for _, sub := range s.routes {
		if sub.Subrouter {
			subs++
		}
	}
	return subs
}
func (s *Server) Routes() int {
	subs := 0
	for _, sub := range s.routes {
		if !sub.Subrouter && sub.Handler != nil {
			subs++
		}
	}
	return subs
}
