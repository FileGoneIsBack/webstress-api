package adminapi

import (
	"api/core/database"
	"api/core/master/sessions"
	"api/core/models/server"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"fmt"
	"strconv"
)
// Custom type for time to handle Unix timestamp in milliseconds
type CustomTime struct {
	time.Time
}

func init() {
	Route.NewSub(server.NewRoute("/invites", func(w http.ResponseWriter, r *http.Request) {
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
			PostAuthTokens(w, r)
		case "get":
			GetAuthTokens(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		}
	}))
}

func PostAuthTokens(w http.ResponseWriter, r *http.Request) {
    type request struct {
        Token  int       `json:"token"`
        Expire time.Time `json:"expire"`
    }
    var req request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    if req.Token <= 0 {
        http.Error(w, "Invalid token", http.StatusBadRequest)
        return
    }
    if req.Expire.Before(time.Now()) {
        http.Error(w, "Expiration time cannot be in the past", http.StatusBadRequest)
        return
    }
    if err := database.Container.AddInvite(req.Token, req.Expire); err != nil {
        http.Error(w, "Failed to add token to the database", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Token Added!"))
}


func GetAuthTokens(w http.ResponseWriter, r *http.Request) {
invites, err := database.Container.GetAllInvites()
if err != nil {
    // Send the error message to the client
    http.Error(w, fmt.Sprintf("Failed to retrieve blacklists: %v", err), http.StatusInternalServerError)
    return
}

	invitesJSON, err := json.Marshal(invites)
	if err != nil {
		http.Error(w, "Failed to marshal blacklists into JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(invitesJSON)
}

func init() {
	Route.NewSub(server.NewRoute("/AddToken", func(w http.ResponseWriter, r *http.Request) {
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

			type request struct {
				Expire CustomTime `json:"exp"`  // Use CustomTime type to handle the Unix timestamp
				Token  string        `json:"token"`
			}

			var req request
			// Decode the incoming JSON request body
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				logger.Println("Error decoding JSON request:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Remove the host from the blacklist in the database
			token, _ := strconv.Atoi(req.Token)
			if err := database.Container.AddInvite(token, req.Expire.Time); err != nil {
				http.Error(w, "Failed to add token to the invite manager", http.StatusInternalServerError)
				return
			}

			// Respond with success status
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Host added from the token successfully"))

		default:
			// Handle unsupported HTTP methods
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		}
	}))
}

// UnmarshalJSON to parse Unix timestamp in milliseconds
func (c *CustomTime) UnmarshalJSON(data []byte) error {
	// Remove quotes from the string (if it's a string) and parse the timestamp as an integer
	timestampStr := strings.Trim(string(data), `"`)
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return fmt.Errorf("unable to parse timestamp: %v", err)
	}

	// Convert from milliseconds to seconds
	c.Time = time.Unix(timestamp/1000, (timestamp%1000)*1000000)
	return nil
}
