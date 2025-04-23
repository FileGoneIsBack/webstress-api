package adminapi

import (
	"api/core/database/site"
	"api/core/database/atks"
	"api/core/database/users"
	"api/core/master/sessions"
	"api/core/models/functions"
	"api/core/models/plans"
	"api/core/models/server"
	"api/core/models/servers"
	"net/http"
	"strings"
)

func init() {
	Route.NewSub(server.NewRoute("/data", func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Method) == "post" {
			ok, session := sessions.IsLoggedIn(w, r)
			if !ok {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}
			if !users.HasPermission(session.User, "admin") {
				http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
				return
			}
			type plan struct {
				Name string `json:"name"`
			}
			type user struct {
				Username string `json:"username"`
				ID       int    `json:"id"`
			}
			type serverInfo struct {
				Name    string   `json:"name"`
				Methods []string `json:"methods"`
			}
			type Data struct {
				UserCount          int            `json:"userCount"`
				AttackCount        int            `json:"dailyAttackCount"`
				RunningAttackCount int            `json:"runningAttackCount"`
				ProfitCount        int            `json:"profitCount"`
				Users              []*user        `json:"users"`
				Plans              []*plan        `json:"plans"`
				Servers            []serverInfo   `json:"servers"`  // New field for methods per server
			}
			d := new(Data)
			d.UserCount = users.Container.Users()
			d.AttackCount = atks.Container.DailyAttacks()
			d.RunningAttackCount = atks.Container.GlobalRunning()
			d.ProfitCount = site.Container.Sales()
			d.Users = func() []*user {
				users, err := users.Container.GetUsers()
				if err != nil {
					return nil
				}
				var u []*user = make([]*user, 0)
				for _, us := range users {
					u = append(u, &user{Username: us.Username, ID: us.ID})
				}
				return u
			}()
			d.Plans = func() []*plan {
				var p []*plan = make([]*plan, 0)
				for name := range plans.Plans {
					p = append(p, &plan{Name: name})
				}
				return p
			}()

			// Fetch server methods
			d.Servers = func() []serverInfo {
				var serversInfo []serverInfo
				for _, srv := range servers.Servers {  // Assuming 'server.Servers' contains all the servers
					serversInfo = append(serversInfo, serverInfo{
						Name:    srv.Name,
						Methods: srv.Methods,
					})
				}
				return serversInfo
			}()

			functions.WriteJson(w, d)
		} else {
			w.Write([]byte("404 page not found"))
			w.WriteHeader(404)
		}
	}))
}
