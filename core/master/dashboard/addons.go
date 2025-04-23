package dashboard

import (
	"api/core/database/users"
	"api/core/database/atks"
	"api/core/database/site"
	"api/core/master/sessions"
	"api/core/models"
	"api/core/models/apis"
	"api/core/models/functions"
	"api/core/models/server"
	"api/core/models/servers"
	"net/http"
)

func init() {
	Route.NewSub(server.NewRoute("/addons", func(w http.ResponseWriter, r *http.Request) {
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
		functions.Render(Page{
			Name:         models.Config.Name,
			Title:        "Addons",
			Vers:         models.Config.Vers,
			ServersCount: len(servers.Servers) + len(apis.Apis),
			Ongoing:      atks.Container.GlobalRunning(),
			Slots:        servers.Slots()[0],
			Users:         users.Container.Users() + site.Site.FakeUsers,
			Remotes:      servers.Servers,
			Session:      user,
		}, w, r, "user", "addons.html")
	}))
}
