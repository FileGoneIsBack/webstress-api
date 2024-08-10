package ticketapi

import (
    "encoding/json"
    "net/http"
    "regexp"
    "strings"

    "api/core/database"
    "api/core/models/server"
    "api/core/master/sessions"
    "github.com/microcosm-cc/bluemonday"
)

func sanitizeInput(input string) string {
    // Use bluemonday to sanitize HTML content
    p := bluemonday.UGCPolicy()
    return p.Sanitize(input)
}

func validateMessage(message string) bool {
    // Example regex to allow only printable characters and spaces
    re := regexp.MustCompile(`^[\p{L}\p{M}\p{N}\p{P}\p{Zs}]+$`)
    return re.MatchString(message)
}

func init() {
    Route.NewSub(server.NewRoute("/update", func(w http.ResponseWriter, r *http.Request) {
        type Status struct {
            Status  string `json:"status"`
            Message string `json:"message"`
        }

        switch strings.ToLower(r.Method) {
        case "post":
            ok, user := sessions.IsLoggedIn(w, r)
            if !ok {
                http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
                return
            }

            var updateticket Ticket
            err := json.NewDecoder(r.Body).Decode(&updateticket)
            if err != nil {
                json.NewEncoder(w).Encode(&Status{Status: "error", Message: "failed to decode request body"})
                return
            }

            // Sanitize and validate the ticket message
            updateticket.Message = sanitizeInput(updateticket.Message)
            if updateticket.Message == "" || !validateMessage(updateticket.Message) {
                json.NewEncoder(w).Encode(&Status{Status: "error", Message: "invalid message"})
                return
            }

            // Save the ticket information in the database
            err = database.Container.UpdateMessage(updateticket.TicketID, user.ID, updateticket.Message)
            if err != nil {
                json.NewEncoder(w).Encode(&Status{Status: "error", Message: err.Error()})
                return
            }

            json.NewEncoder(w).Encode(&Status{Status: "success", Message: "ticket submitted successfully"})
        }
    }))
}
