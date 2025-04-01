package functions

import (
	"api/core/models/log"
	"net/http"
	"path/filepath"
	"text/template"
	"api/core/master/sessions"
	"api/core/database"
	"api/core/models"
	"api/core/models/apis"
	"api/core/models/servers"
)

func Render(v interface{}, w http.ResponseWriter, r *http.Request, file ...string) {
    // Get the flash messages from the session
    flashMessages := sessions.GetFlash(w, r)
	_, userSession := sessions.IsLoggedIn(w, r)
    // Ensure data is a map if v is nil, or cast it to a map if it's already a map
    var data map[string]interface{}
    if v == nil {
        data = make(map[string]interface{})
    } else {
        // If v is already a map, use it, otherwise create a new map
        if m, ok := v.(map[string]interface{}); ok {
            data = m
        } else {
            data = make(map[string]interface{})
            data["data"] = v
        }
    }

    // Add the flash messages to the data
	data = map[string]interface{}{
		"Name":          models.Config.Name,
		"Title":         "Dashboard",
		"Vers":          models.Config.Vers,
		"ServersCount":  len(servers.Servers) + len(apis.Apis),
		"Ongoing":       database.Container.GlobalRunning(),
		"Slots":         servers.Slots()[0],
		"Users":         database.Container.Users() + models.Config.Fake.Users,
		"Remotes":       servers.Servers,
		"Session":       userSession, // Retrieve user session
		"FlashMessages": flashMessages,
	}

    // Parse the main template
    t, err := template.ParseFiles("assets/html/" + filepath.Join(file...))
    if err != nil {
        log.Println(err)
        return
    }

    t, err = t.ParseFiles("assets/html/nav.html")
    if err != nil {
        log.Println(err)
        return
    }
    t, err = t.ParseFiles("core/footer.html")
    if err != nil {
        log.Println(err)
        return
    }
    t, err = t.ParseFiles("assets/html/construction.html")
    if err != nil {
        log.Println(err)
        return
    }

    // Execute the main template
    err = t.Execute(w, data)
    if err != nil {
        log.Println(err)
        return
    }
}
