package authentication

import (
	"api/core/master/internal"
	"api/core/models/server"
	"net/http"
	"strings"
)


func init() {
	Route.NewSub(server.NewRoute("/signup", func(w http.ResponseWriter, r *http.Request) {
		switch strings.ToLower(r.Method) {
		case "get":

		case "post":
			internal.Signup(w, r)
		}
	}))
}
