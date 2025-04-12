package api

import (
	accountapi "api/core/master/api/account"
	adminapi "api/core/master/api/admin"
	attackapi "api/core/master/api/attacks"
	dashboardapi "api/core/master/api/dashboard"
	Managerapi "api/core/master/api/manager"

	panelapi "api/core/master/api/panel"
	paymentsapi "api/core/master/api/payments"
	ticketapi "api/core/master/api/tickets"
	"api/core/models/server"
)

var (
	Route *server.Route = server.NewSubRouter("/api")
)

func init() {
	Route.NewSubs(
		paymentsapi.Route,
		panelapi.Route,
		dashboardapi.Route,
		adminapi.Route,
		ticketapi.Route,
		attackapi.Route,
		accountapi.Route,
		Managerapi.Route, 
	)
}
