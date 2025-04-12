package Managerapi

import "api/core/models/server"

var (
	Route *server.Route = server.NewSubRouter("/manager")
)
