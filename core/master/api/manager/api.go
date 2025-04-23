package Managerapi

import (
	"api/core/database/users"
	"api/core/database/site"
	"api/core/master/sessions"
	"api/core/models/server"
	"api/core/models"
	"encoding/json"
	"net/http"
	"strings"
	"fmt"
)


func init() {
	Route.NewSub(server.NewRoute("/api", func(w http.ResponseWriter, r *http.Request) {
		ok, session := sessions.IsLoggedIn(w, r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		switch strings.ToLower(r.Method) {
		case "post":
			newKey := users.GenerateUserKey(models.Config.Name)
			err := users.Container.UpdateUserKey(session.User.Username, []byte(newKey))
			if err != nil {
				http.Error(w, "Failed to update user key", http.StatusInternalServerError)
				return
			}
			// Respond with the new key
			response := map[string]string{
				"message": "User key updated successfully",
				"new_key": string(newKey),
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		case "get":
			user, _ := users.Container.GetUser(session.User.Username)
			key, err := site.Container.GetApiKey(user.Username)
			if err != nil {
				fmt.Println("Error decoding base64:", err)
				return
			}

			response := map[string]string{
				"key": key, 
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
		}
	}))
}

