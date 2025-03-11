package accountapi

import (
	"os"
	"log"
	"api/core/models/server"
)

var (
	Route *server.Route = server.NewSubRouter("/account")
	logger = log.New(os.Stderr, "[Admin] ", log.Ltime|log.Lshortfile)
)
