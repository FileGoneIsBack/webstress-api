package admin

import (
	"api/core/database/users"
	"api/core/master/sessions"
	"api/core/models"
	"api/core/models/functions"
	"api/core/models/server"
	"api/core/models/servers"
	"api/core/models/floods"

	"net/http"
)

func init() {
	Route.NewSub(server.NewRoute("/methods", func(w http.ResponseWriter, r *http.Request) {
        type Page struct {
            Name, Title, Vers            string
            ServersCount, Ongoing, Slots int
            Users                        int
            TopMethod                    string  // name of most-used
            VipMethods                   int     // # of vip methods
            Usage                        int     // total invocations
            Methods                      int     // # of methods loaded
            Remotes                      map[string]*servers.Server
            *sessions.Session
        }
		top, vipCount, totalUsage, methodCount := floods.Stats()
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
            Title:        "Admin Management",
            Vers:         models.Config.Vers,
            ServersCount: len(servers.Servers),
            TopMethod:    func() string { if top!=nil { return top.Name }; return "" }(),
            VipMethods:   vipCount,
            Usage:        totalUsage,
            Methods:      methodCount,
            Session:      user,
        }, w, r, "admin", "methods.html")
	}))
}
