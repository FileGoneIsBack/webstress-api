package master

import (
	"net/http"
)

// Define the 404 page handler
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "assets/html/404.html")
    
}