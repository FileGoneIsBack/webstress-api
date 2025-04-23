package adminapi

import (
	"api/core/database/site"
	"api/core/database/users"
	"api/core/master/sessions"
	"api/core/models/log"
	"api/core/models/server"
	"encoding/json"
	"net/http"
	"strings"
)

func init() {
	Route.NewSub(server.NewRoute("/news", func(w http.ResponseWriter, r *http.Request) {
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
			var newNews site.News
			if err := json.NewDecoder(r.Body).Decode(&newNews); err != nil {
				log.Println("Error decoding JSON request:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			newNews.From = session.User.Username
			if err := site.Container.NewNews(&newNews); err != nil {
				log.Println("Error adding news to the database:", err)
				http.Error(w, "Failed to add news", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("News added successfully"))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 page not found"))
		}
	}))
}
