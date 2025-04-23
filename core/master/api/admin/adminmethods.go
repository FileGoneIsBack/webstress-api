package adminapi

import (
	"api/core/database/atks"
	"api/core/models/floods"
	"api/core/database/users"
	"api/core/master/sessions"
	"api/core/models/server"

	"encoding/json"
	"net/http"
	"fmt"
)


func init() {
	Route.NewSub(server.NewRoute("/method", func(w http.ResponseWriter, r *http.Request) {
		ok, session := sessions.IsLoggedIn(w, r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		if !users.HasPermission(session.User, "admin") {
			http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
			return
		}
	
		switch r.Method {
		case http.MethodGet:
			// load from DB into in‑memory cache
			if err := atks.Container.LoadFloodMethods(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// return the map[string]*floods.Method as JSON
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(floods.Methods)
	
		case http.MethodPost:
			// generic request shape
			var req struct {
				Action string        `json:"action"`
				Method floods.Method `json:"method"`
				Name   string        `json:"name"`   
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid payload: "+err.Error(), http.StatusBadRequest)
				return
			}
			
			var err error
			switch req.Action {
			case "add":
				err = atks.Container.AddFloodMethod(&req.Method)
			case "update":
				err = atks.Container.UpdateFloodMethod(&req.Method, req.Name)
			case "delete":
				err = atks.Container.DeleteFloodMethod(req.Name)
			default:
				http.Error(w, "unknown action", http.StatusBadRequest)
				return
			}
	
			if err != nil {
				http.Error(w, fmt.Sprintf("%s failed: %v", req.Action, err), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"ok"}`))
	
		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))
}