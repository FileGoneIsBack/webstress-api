package attackapi

import (
	"api/core/database"
	"api/core/master/sessions"
	"api/core/models/functions"
	"api/core/models/server"
	"api/core/models/servers"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func init() {
	Route.NewSub(server.NewRoute("/stop/{id}", func(w http.ResponseWriter, r *http.Request) {
		type status struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}

		switch strings.ToLower(r.Method) {
		case "get":
			key, ok := functions.GetKey(w, r)
			if !ok {
				return
			}
			
			vars := mux.Vars(r)
			attackID := vars["id"]
			servers.Stop(key.ID, attackID)
			if err := database.Container.Stop(key, attackID); err != nil {
				json.NewEncoder(w).Encode(&status{
					Status:  "error",
					Message: "failed to stop attack!",
				})
				return
			}

			json.NewEncoder(w).Encode(&status{
				Status: "success",
			})

		case "post":
			ok, user := sessions.IsLoggedIn(w, r)
			if !ok {
				return
			}

			vars := mux.Vars(r)
			attackID := vars["id"]
			servers.Stop(user.ID, vars["id"])
			if err := database.Container.Stop(user.User, attackID); err != nil {
				json.NewEncoder(w).Encode(&status{
					Status:  "error",
					Message: "failed to stop attack!",
				})
				return
			}

			json.NewEncoder(w).Encode(&status{
				Status: "success",
			})
		}
	}))
}
