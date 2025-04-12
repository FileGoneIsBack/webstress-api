package internal

import (
	"api/core/database"
	"api/core/master/sessions"
	"net/http"
	"time"
	"api/core/models/ranks"
	"github.com/google/uuid"
)

func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}
	
	username := r.FormValue("login-username")
	if username == "" {
		renderErrorPage(w, r, "Missing fields", "error")
		return
	}

	password := r.FormValue("login-password")
	if password == "" {
		renderErrorPage(w, r, "Missing fields", "error")
		return
	}

	user, err := database.Container.GetUser(r.Form["login-username"][0])
	if err != nil {
		renderErrorPage(w, r, "Invalid username or password", "error")
		return
	}

	if user == nil {
		renderErrorPage(w, r, "Invalid username or password", "error")
		return
	}
	if !user.IsKey([]byte(r.Form["login-password"][0])) {
		renderErrorPage(w, r, "Invalid username or password", "error")
		return
	}
	expiryTime := time.Unix(user.Expiry, 0)
	if expiryTime.Before(time.Now()) {
		user = &database.User{
		Username: user.Username,
		Membership: "expired",
		Expiry: time.Now().Add(31 * 24 * time.Hour).Unix(),
		Concurrents: 0,
		Servers: 0,
		Duration: 0,
		Ranks: []*ranks.Rank{
			ranks.GetRole("member", true),
		},
		}
		err := database.Container.UpdateUser(user)
		if err != nil {
			renderErrorPage(w, r, "Error updating account! contact staff or try again", "error")
			return
		}
	}


	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(30 * time.Minute)
	if _, remember := r.Form["remember-me"]; remember {
		expiresAt = time.Now().Add(48 * time.Hour)
	}

	sessions.Sessions[sessionToken] = sessions.Session{
		User:   user,
		Expiry: expiresAt,
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session-token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}
