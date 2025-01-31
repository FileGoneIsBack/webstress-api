package internal

import (
	"api/core/database"
	"api/core/master/sessions"
	"api/core/models"
	"api/core/models/functions"
	"html/template"
	"net/http"
	"time"
	"api/core/models/ranks"
	"github.com/google/uuid"
)

func Login(w http.ResponseWriter, r *http.Request) {
	type Page struct {
		Name   string
		Title  string
		Script template.HTML
	}
	err := r.ParseForm()
	if err != nil {
		return
	}
	
	username := r.FormValue("login-username")
	if username == "" {
		functions.Render(Page{
			Name:  models.Config.Name,
			Title: "Login",
			Script: template.HTML(functions.Toast(functions.Toastr{
				Icon:  "error",
				Title: "Error!",
				Text:  "Username is required.",
			})),
		}, w, "login", "login.html")
		return
	}

	password := r.FormValue("login-password")
	if password == "" {
		functions.Render(Page{
			Name:  models.Config.Name,
			Title: "Login",
			Script: template.HTML(functions.Toast(functions.Toastr{
				Icon:  "error",
				Title: "Error!",
				Text:  "Password is required.",
			})),
		}, w, "login", "login.html")
		return
	}

	user, err := database.Container.GetUser(r.Form["login-username"][0])
	if err != nil {
		functions.Render(Page{
			Name:  models.Config.Name,
			Title: "Login",
			Script: template.HTML(functions.Toast(functions.Toastr{
				Icon:  "error",
				Title: "Error!",
				Text:  "Invalid credentials.",
			})),
		}, w, "login", "login.html")
		return
	}

	if user == nil {
		functions.Render(Page{
			Name:  models.Config.Name,
			Title: "Login",
			Script: template.HTML(functions.Toast(functions.Toastr{
				Icon:  "error",
				Title: "Error!",
				Text:  "Invalid credentials.",
			})),
		}, w, "login", "login.html")
		return
	}
	if !user.IsKey([]byte(r.Form["login-password"][0])) {
		functions.Render(Page{
			Name:  models.Config.Name,
			Title: "Login",
			Script: template.HTML(functions.Toast(functions.Toastr{
				Icon:  "error",
				Title: "Error!",
				Text:  "Invalid credentials.",
			})),
		}, w, "login", "login.html")
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
			functions.Render(Page{
				Name:  models.Config.Name,
				Title: "Login",
				Script: template.HTML(functions.Toast(functions.Toastr{
					Icon:  "error",
					Title: "Error!",
					Text:  "There was an error updating your account status.",
				}))},
				w, "login", "login.html")
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
