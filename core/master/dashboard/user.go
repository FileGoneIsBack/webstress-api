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
    Route.NewSub(server.NewRoute("/dashboard", func(w http.ResponseWriter, r *http.Request) {
        type Page struct {
            Name, Title, Vers       string
            ServersCount, Ongoing, Slots int
            Users                   int
            Remotes                 map[string]*servers.Server
            *sessions.Session
            FlashMessages           []sessions.FlashMessage        
        }
        ok, user := sessions.IsLoggedIn(w, r)
        if !ok {
            http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
            return
        }
        //system noti
        sessions.SetFlash(w, r, "Your changes have been saved successfully.", "System")
        
        flashMessages := sessions.GetFlash(w, r)
        functions.Render(Page{
            Name:    models.Config.Name,
            Title:   "Dashboard",
            Vers:    models.Config.Vers,
            ServersCount: len(servers.Servers) + len(apis.Apis),
            Ongoing: database.Container.GlobalRunning(),
            Slots:   servers.Slots()[0],
            Users:   database.Container.Users() + models.Config.Fake.Users,
            Remotes: servers.Servers,
            Session: user,
            FlashMessages: flashMessages,
        }, w, r, "user", "dash.html")
    }))
}
