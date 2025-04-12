package api

import (
	//"encoding/json"
	"bot/models"
	//"fmt"
	//"io/ioutil"
	//"log"
	"net/http"
	"strings"
	//"bot/database"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		/*if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Check for Auth
		authHeader := r.Header.Get("Authorization")
		if !isValidToken(authHeader) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Print("Failed invalid key or someone is attempting to hack the bot!")
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			log.Printf("Received body: %s", string(body))
			return
		}
		log.Printf("Received body: %s", string(body))
		var requestUser struct {
			Username   string `json:"username"`
			Telegram   string `json:"tgusername"`
			Membership string `json:"membership"`
			Balance    int    `json:"balance"`
		}

		err = json.Unmarshal(body, &requestUser)
		if err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
			return
		}

		user := &database.User{
			Username:   requestUser.Username,
			Telegram:   requestUser.Telegram,
			Membership: requestUser.Membership,
			Balance:    requestUser.Balance,
		}
		err = database.Container.AddUser(user)
		if err != nil {
			http.Error(w, "Failed to add user to the database", http.StatusInternalServerError)
			log.Printf("Error adding user: %v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "success", "message": "User added to database"}`)
		*/
	})
}

func isValidToken(authHeader string) bool {
	storedToken := models.Config.Auth
	return strings.HasPrefix(authHeader, "Bearer "+storedToken)
}
