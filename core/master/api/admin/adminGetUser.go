package adminapi

import (
	"api/core/database/site"
	"api/core/database/users"
	"api/core/master/sessions"
	"api/core/models/server"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)


func init() {
	Route.NewSub(server.NewRoute("/GetUser", func(w http.ResponseWriter, r *http.Request) {
		ok, session := sessions.IsLoggedIn(w, r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		if !users.HasPermission(session.User, "admin") {
			http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
			return
		}
	
		switch strings.ToLower(r.Method) {
		case "get":
			ticketCookie, err := r.Cookie("ticket_id")
			if err != nil {
				http.Error(w, "Missing ticket_id cookie", http.StatusBadRequest)
				return
			}
			ticketIDStr := ticketCookie.Value
			ticketID, _ := strconv.Atoi(ticketIDStr)
			log.Printf("TESTING STRING: %s", ticketIDStr)
			user, err := site.Container.GetUserForTicket(ticketID)
			if err != nil {
				http.Error(w, "Error fetching user: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if user == nil {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
	
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
	
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		}
	}))
}