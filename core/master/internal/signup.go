package internal

import (
	"api/core/database"
	"api/core/database/users"
	"api/core/database/site"
	"api/core/master/sessions"
	"api/core/models/antiflood"
	"api/core/models/ranks"
	"errors"
	"fmt"
	"html/template"
	"io"
	"regexp"
	"net/url"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)
type Page struct {
	Name   string
	Title  string
	Script template.HTML
}

var (
	SignupLimiter = antiflood.NewRateLimiter(
		1,
		120*time.Minute,
		antiflood.WithKeyByRealIP(),
	)
)

func GetBodyData(data string) string {
	fmt.Println(data)
	return strings.Split(data, "=")[1]
}

func KeyByRealIP(r *http.Request) string {
	var ip string

	if tcip := r.Header.Get("True-Client-IP"); tcip != "" {
		ip = tcip
	} else if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		ip = xrip
	} else if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		i := strings.Index(xff, ", ")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	} else if ccip := r.Header.Get("CF-Connecting-IP"); ccip != "" {
		ip = ccip
	} else {
		var err error
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
	}

	return ip
}

func Signup(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	data := strings.Split(string(body), "&")
	fmt.Println(data)

	if len(data) < 5 {
		renderErrorPage(w, r, "Please fill out all fields!", "error")
		return
	}

	username := GetBodyData(data[0])
	password, err := url.QueryUnescape(GetBodyData(data[1]))
	if err != nil {
		fmt.Println("Error during URL decoding:", err)
	}
	fmt.Println("Decoded password:", password)
	cpassword, _ := url.QueryUnescape(GetBodyData(data[2]))
	auth	 := GetBodyData(data[3])
	tos := GetBodyData(data[4])

	// Validate form data
	if err := validateSignupData(username, password, cpassword, tos); err != nil {
		renderErrorPage(w, r, err.Error(), "error")
		return
	}

	//check if key is valid
	// Check if user already exists
	user, err := users.Container.GetUser(username)
	if err != nil && !errors.Is(err, users.ErrUserNotFound) {
		renderDatabaseErrorPage(w, r, "Error retrieving user from database.")
		return
	}
	if user != nil {
		renderErrorPage(w, r, "User already exists.", "error")
		return
	}
	tele, err := users.Container.CheckInvite(auth)
	if err != nil {
		renderErrorPage(w, r, err.Error(), "error")
		return
	}
	// Create new user
	user = &database.User{
		Username:    username,
		Key:         []byte(password),
		Membership: "Free",
		Ranks: []*ranks.Rank{
			ranks.GetRole("member", true),
		},
		Tele: tele,
	}
	if site.Site.FreeUser {
		user.Concurrents = site.Site.FreeUserCons
		user.Servers = 0
		user.Duration = site.Site.FreeUserTime
		user.Expiry = time.Now().Add(31 * 24 * time.Hour).Unix()
	}
	err = users.Container.NewUser(user)
	if err != nil {
		renderDatabaseErrorPage(w, r, "Error creating new user in database.")
		return
	}

	// Set session and redirect to dashboard
	user, _ = users.Container.GetUser(username)
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(30 * time.Minute)
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

func renderErrorPage(w http.ResponseWriter, r *http.Request, errorMessage string, errortype string) {
	// URL encode the error message to pass it in the query string
	encodedMessage := url.QueryEscape(errorMessage + errortype)

	// Redirect the user back to the form with the error message
	http.Redirect(w, r, "/login?error=" + encodedMessage, http.StatusFound)
}

func renderDatabaseErrorPage(w http.ResponseWriter, r *http.Request, errorMessage string) {
	renderErrorPage(w, r, "Database error occurred: "+errorMessage, "error")
}


func validateSignupData(username, password, cpassword, tos string) error {
	if len(username) < 4 {
		return errors.New("username must be at least 4 characters")
	}
	usernameRe := regexp.MustCompile(`^[a-zA-Z0-9!@#$_-]+$`)
	if !usernameRe.MatchString(username) {
		return fmt.Errorf("username can only contain letters & numbers")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 6 characters")
	}
	passwordRe := regexp.MustCompile(`^[a-zA-Z0-9!@#$_-]+$`)
	if !passwordRe.MatchString(password) {
		return fmt.Errorf("password can only contain letters, numbers, and !@#$_-")
	}
	if password != cpassword {
		return errors.New("passwords do not match")
	}
	if tos != "on" {
		return errors.New("you must agree to the terms of service")
	}

	return nil
}