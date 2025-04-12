package functions

import (
	"api/core/models/log"
	"net/http"
    "reflect"
	"path/filepath"
	"text/template"
	"api/core/master/sessions"
	"api/core/database"
	"api/core/models"
	"api/core/models/apis"
	"api/core/models/servers"
)

func Render(v interface{}, w http.ResponseWriter, r *http.Request, file ...string) {
    flashMessages := sessions.GetFlash(w, r)
    _, userSession := sessions.IsLoggedIn(w, r)

    data := make(map[string]interface{})

    // Convert struct or use map
    if v != nil {
        switch x := v.(type) {
        case map[string]interface{}:
            for k, val := range x {
                data[k] = val
            }
        default:
            for k, val := range structToMap(v) {
                data[k] = val
            }
        }
    }

    // Add global data (without overwriting existing keys)
    global := map[string]interface{}{
        "Name":          models.Config.Name,
        "Title":         "test",
        "Vers":          models.Config.Vers,
        "ServersCount":  len(servers.Servers) + len(apis.Apis),
        "Ongoing":       database.Container.GlobalRunning(),
        "Slots":         servers.Slots()[0],
        "Users":         database.Container.Users() + models.Config.Fake.Users,
        "Remotes":       servers.Servers,
        "Session":       userSession,
        "FlashMessages": flashMessages,
    }

    for k, val := range global {
        if _, exists := data[k]; !exists {
            data[k] = val
        }
    }

    // Parse templates
    t, err := template.ParseFiles("assets/html/" + filepath.Join(file...))
    if err != nil {
        log.Println(err)
        return
    }
    t, err = t.ParseFiles("assets/html/nav.html", "core/footer.html", "assets/html/construction.html")
    if err != nil {
        log.Println(err)
        return
    }

    // Render
    err = t.Execute(w, data)
    if err != nil {
        log.Println(err)
    }
}


func structToMap(v interface{}) map[string]interface{} {
    result := make(map[string]interface{})
    val := reflect.ValueOf(v)
    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }
    typ := val.Type()

    for i := 0; i < val.NumField(); i++ {
        field := typ.Field(i)
        if field.PkgPath != "" { // skip unexported fields
            continue
        }
        result[field.Name] = val.Field(i).Interface()
    }

    return result
}