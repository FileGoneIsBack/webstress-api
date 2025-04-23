package admin

import (
	"api/core/database/site"
	"api/core/database/users"
	"api/core/database/atks"
	"api/core/master/sessions"
	"api/core/models"
	"api/core/models/apis"
	"api/core/models/functions"
	"api/core/models/server"
	"api/core/models/servers"
	"net/http"
)

func init() {
	Route.NewSub(server.NewRoute("/servers", func(w http.ResponseWriter, r *http.Request) {
		type Page struct {
			Name, Title, Vers            string
			ServersCount, Ongoing, Slots int
			Users                        int
			Remotes                      map[string]*servers.Server
			*sessions.Session
		}
		ok, user := sessions.IsLoggedIn(w, r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		if !users.HasPermission(user.User, "admin") {
			http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
			return
		}
		functions.Render(Page{
			Name:         models.Config.Name,
			Title:        "Servers",
			Vers:         models.Config.Vers,
			ServersCount: len(servers.Servers) + len(apis.Apis),
			Ongoing:      atks.Container.GlobalRunning(),
			Slots:        servers.Slots()[0],
			Users:         users.Container.Users() + site.Site.FakeUsers,
			Remotes:      servers.Servers,
			Session:      user,
		}, w, r, "admin", "servers.html")
	}))
}
