package server

import (
	"api/core/models"
	_ "api/core/models/antiflood"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/http2"
)
var logger = log.New(os.Stdout, "[TLS/Server] ", log.LstdFlags)
type Server struct {
	server *http.Server
	router *mux.Router
	routes map[string]*Route
	logger *log.Logger
	config *Config
}

func NewServer(config *Config) *Server {
	s := &Server{
		server: &http.Server{
			Addr:         config.Addr,
			Handler:      nil,
			ErrorLog:     logger,
			WriteTimeout: 30 * time.Second, // Increased from 15 seconds
			ReadTimeout:  30 * time.Second, // Increased from 15 seconds
		},
		router: mux.NewRouter(),
		routes: make(map[string]*Route),
		logger: log.New(os.Stderr, "[server] ", log.Ltime|log.Lshortfile),
		config: config,
	}

	// Configure HTTP/2
	http2.ConfigureServer(s.server, &http2.Server{
		MaxConcurrentStreams: 50, // Adjust based on your needs
		IdleTimeout: 30 * time.Minute,
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
		logger.Print("Server is running on HTTPS on " + s.server.Addr)
		return s.server.ListenAndServeTLS(cert, key)
	} else {
		logger.Print("Server is running on HTTP on " + s.server.Addr)
	}

	s.logger.Println("listening with " + fmt.Sprint(s.Subrouters()) + " subrouters and " + fmt.Sprint(s.Routes()) + " routes.")

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
