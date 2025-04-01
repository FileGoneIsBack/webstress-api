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
	Route.NewSub(server.NewRoute("/account", func(w http.ResponseWriter, r *http.Request) {
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
		err := sessions.SetFlash(w, r, "Your changes have been saved successfully.")
		if err != nil {
			http.Error(w, "Error setting flash message: "+err.Error(), http.StatusInternalServerError)
			return
		}
		functions.Render(Page{
			Name:         models.Config.Name,
			Title:        "Profile",
			Vers:         models.Config.Vers,
			ServersCount: len(servers.Servers) + len(apis.Apis),
			Ongoing:      database.Container.GlobalRunning(),
			Slots:        servers.Slots()[0],
			Users:        database.Container.Users() + models.Config.Fake.Users,
			Remotes:      servers.Servers,
			Session:      user,
		}, w, r, "user", "account.html")
	}))
}
