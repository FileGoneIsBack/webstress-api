package accountapi

import (
	//"api/core/database"
	"api/core/master/sessions"
	"api/core/models"
	"api/core/models/server"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

func init() {
	Route.NewSub(server.NewRoute("/update", func(w http.ResponseWriter, r *http.Request) {
		ok, _ := sessions.IsLoggedIn(w, r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		switch strings.ToLower(r.Method) {
		case "post":
			handlePostRequest(w, r)
		case "get":
			//handleGetRequest(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		}
	}))
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	ok, user := sessions.IsLoggedIn(w, r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	var request struct {
		TGusername string `json:"telegram"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		logger.Print("Failed to decode JSON body") 
		return
	}

	payload := map[string]interface{}{
		"username":   user.Username,
		"tgusername": request.TGusername,
		"membership": user.Membership,
		"balance":    user.Balance,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		logger.Print("Failed to marshal JSON payload") 
		return
	}

	apiURL := models.Config.Bot.URL 
	resp, err := sendAPIRequest(apiURL, jsonData)
	if err != nil {
		http.Error(w, "Failed to send request to API", http.StatusInternalServerError)
		logger.Print("Failed to send API request: " + err.Error()) 
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		w.WriteHeader(http.StatusOK)
		logger.Printf("User %s passed telegram auth with @%s", user.Username, request.TGusername)
	} else {
	}
}

func sendAPIRequest(url string, jsonData []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	storedToken := models.Config.Bot.Auth
	req.Header.Set("Authorization", "Bearer "+storedToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
