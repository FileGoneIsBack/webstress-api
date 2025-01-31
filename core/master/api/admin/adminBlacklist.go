package adminapi

import (
	"api/core/database"
	"api/core/master/sessions"
	"api/core/models/server"
	"encoding/json"
	"net/http"
	"strings"
)

func init() {
	Route.NewSub(server.NewRoute("/blacklist", func(w http.ResponseWriter, r *http.Request) {
		ok, session := sessions.IsLoggedIn(w, r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		if !session.HasPermission("admin") {
			http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
			return
		}
		switch strings.ToLower(r.Method) {
		case "post":
			handlePostRequest(w, r)
		case "get":
			handleGetRequest(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		}
	}))
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	ok, session := sessions.IsLoggedIn(w, r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	if !session.HasPermission("admin") {
		http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
		return
	}

	var request struct {
		Host string `json:"host"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Println("Error decoding JSON request:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := database.Container.NewBlacklist(request.Host); err != nil {
		logger.Println("Error adding host to the blacklist:", err)
		http.Error(w, "Failed to add host to the blacklist", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Host added to the blacklist successfully"))
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	blacklists, err := database.Container.GetAllBlacklists()
	if err != nil {
		logger.Println("Error retrieving blacklists from the database:", err)
		http.Error(w, "Failed to retrieve blacklists", http.StatusInternalServerError)
		return
	}

	blacklistJSON, err := json.Marshal(blacklists)
	if err != nil {
		logger.Println("Error marshaling blacklists into JSON:", err)
		http.Error(w, "Failed to marshal blacklists into JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(blacklistJSON)
}

func init() {
	Route.NewSub(server.NewRoute("/removeBlacklist", func(w http.ResponseWriter, r *http.Request) {
		switch strings.ToLower(r.Method) {
		case "post":
			ok, session := sessions.IsLoggedIn(w, r)
			if !ok {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}
			if !session.HasPermission("admin") {
				http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
				return
			}

			var request struct {
				Host string `json:"host"`
			}
			// Decode the incoming JSON request body
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				logger.Println("Error decoding JSON request:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Remove the host from the blacklist in the database
			if err := database.Container.RemoveBlacklist(request.Host); err != nil {
				logger.Println("Error removing host from the blacklist:", err)
				http.Error(w, "Failed to remove host from the blacklist", http.StatusInternalServerError)
				return
			}

			// Respond with success status
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Host removed from the blacklist successfully"))
		default:
			// Handle unsupported HTTP methods
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		}
	}))
}
