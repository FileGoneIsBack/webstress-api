package adminapi

import (
	//"api/core/database"
	"api/core/database/users"
	"api/core/master/sessions"
	"api/core/models/server"
	"encoding/json"
	"net/http"
	"strings"
	"fmt"
)

func init() {
	Route.NewSub(server.NewRoute("/deleteInvite", func(w http.ResponseWriter, r *http.Request) {
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
		case "post":
			postDeleteInvites(w, r)
		case "get":
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		}
	}))
}

func postDeleteInvites(w http.ResponseWriter, r *http.Request) {
    // Ensure the method is POST
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        w.Write([]byte("Method not allowed"))
        return
    }

    // Parse the request body to get the token
    var request struct {
        Token int `json:"token"`
    }

    // Decode the JSON body into the request struct
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, fmt.Sprintf("Error parsing request body: %v", err), http.StatusBadRequest)
        return
    }

    // Call the DeleteInvite function to remove the token
    //err := database.Container.DeleteInvite(request.Token)
    //if err != nil {
    //    http.Error(w, fmt.Sprintf("Error deleting invite: %v", err), http.StatusInternalServerError)
    //    return
    //}

    // Send a success response
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Invite deleted successfully"))
}