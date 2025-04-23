package adminapi

import (
	"api/core/database"
	"api/core/database/users"
	"api/core/master/sessions"
	"api/core/models/log"
	"api/core/models/server"
	"encoding/json"
	"net/http"
	"strings"
)

func init() {
	Route.NewSub(server.NewRoute("/update-user", func(w http.ResponseWriter, r *http.Request) {
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
			var updatedUser database.User
			if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
				log.Println("Error decoding JSON request:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Update the user in the database
			if err := users.Container.UpdateUser(&updatedUser); err != nil {
				log.Println("Error updating user in the database:", err)
				http.Error(w, "Failed to update user", http.StatusInternalServerError)
				return
			}

			// Respond with success status
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("User updated successfully"))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 page not found"))
		}
	}))
}

func init() {
	Route.NewSub(server.NewRoute("/delete-user", func(w http.ResponseWriter, r *http.Request) {
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
			var deleteUserReq struct {
				Username string `json:"username"`
			}
			if err := json.NewDecoder(r.Body).Decode(&deleteUserReq); err != nil {
				log.Println("Error decoding JSON request:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Fetch the user ID
			userID, err := users.Container.GetUserID(deleteUserReq.Username)
			if err != nil {
				log.Println("Error fetching user ID from the database:", err)
				http.Error(w, "Failed to delete user", http.StatusInternalServerError)
				log.Println(deleteUserReq.Username)
				return
			}

			// Delete the user
			if err := users.Container.DeleteUser(deleteUserReq.Username, userID); err != nil {
				log.Println("Error deleting user from the database:", err)
				http.Error(w, "Failed to delete user", http.StatusInternalServerError)
				return
			}

			// Respond with success status
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("User deleted successfully"))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 page not found"))
		}
	}))
}
