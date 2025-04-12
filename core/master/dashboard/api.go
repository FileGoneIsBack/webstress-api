package dashboard

import (
	"api/core/database"
	"api/core/master/sessions"
	"api/core/models"
	"api/core/models/apis"
	"api/core/models/functions"
	"api/core/models/server"
	"api/core/models/servers"
	"net/http"
)

func init() {
	Route.NewSub(server.NewRoute("/api_manager", func(w http.ResponseWriter, r *http.Request) {
		type Page struct {
			Name, Title, Vers, Domain    string
			ServersCount, Ongoing, Slots int
			Users, Reqs, FailedReqs      int
			Remotes                      map[string]*servers.Server
			SuccessRate 				 int
			*sessions.Session
		}
		ok, user := sessions.IsLoggedIn(w, r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		Reqs, FailedReqs, _ := database.Container.GetReqs(user.Username)
		successRate := 100
		if Reqs > 0 {
			successRate = int(float64(Reqs-FailedReqs) / float64(Reqs) * 100)
		}
		functions.Render(Page{
			Name:         models.Config.Name,
			Title:        "Manager",
			Vers:         models.Config.Vers,
			ServersCount: len(servers.Servers) + len(apis.Apis),
			Ongoing:      database.Container.GlobalRunning(),
			Slots:        servers.Slots()[0],
			Users:        database.Container.Users() + models.Config.Fake.Users,
			Remotes:      servers.Servers,
			Session:      user,
			Domain:       models.Config.Domain,
			Reqs: 		  Reqs,
			FailedReqs:	  FailedReqs,
			SuccessRate:  successRate,
		}, w, r, "api", "api.html")
	}))
}
