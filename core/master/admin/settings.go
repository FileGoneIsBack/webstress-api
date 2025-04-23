package admin

import (
	"api/core/database/site"
	"api/core/database/users"
	"api/core/database/atks"
	"api/core/master/sessions"
	"api/core/models"
	"api/core/models/functions"
	"api/core/models/server"
	"api/core/models/servers"
	"net/http"
)

func init() {
	Route.NewSub(server.NewRoute("/settings", func(w http.ResponseWriter, r *http.Request) {
		type Page struct {
			Name, Title, API, Secure	 string
			FreeUser				     string
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
	    apiStr := "Disabled"
		if site.Site != nil && site.Site.API {
			apiStr = "Active"
		}
		secureStr := "Disabled"
		if models.Config.Secure {
			secureStr = "Active"
		}
		FreeUserStr := "Disabled"
		if site.Site != nil && site.Site.FreeUser {
			FreeUserStr = "Active"
		}
		functions.Render(Page{
			Name:         models.Config.Name,
			Title:        "Settings",
			API:          apiStr,
			FreeUser:	  FreeUserStr,
			Secure:		  secureStr,
			Ongoing:      atks.Container.GlobalRunning(),
			Slots:        servers.Slots()[0],
			Users:         users.Container.Users() + site.Site.FakeUsers,
			Session:      user,
		}, w, r, "admin", "settings.html")
	}))
}
