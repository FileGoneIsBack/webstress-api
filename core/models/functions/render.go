package functions

import (
    "api/core/database/atks"
    "api/core/database/site"
    "api/core/database/users"
    "api/core/master/sessions"
    "api/core/models"
    "api/core/models/apis"
    "api/core/models/log"
    "api/core/models/servers"
    "net/http"
    "path/filepath"
    "reflect"
    "text/template"
)

type RenderData struct {
    Name          string
    Title         string
    Vers          string
    ServersCount  int
    Ongoing       int
    Slots         int
    Users         int
    Remotes       map[string]*servers.Server
    Session       *sessions.Session
    FlashMessages []sessions.FlashMessage
    Payload       map[string]interface{}
}

func Render(v interface{}, w http.ResponseWriter, r *http.Request, file ...string) {
    flashMessages := sessions.GetFlash(w, r)
    _, userSession := sessions.IsLoggedIn(w, r)

    // build your data map as before…
    data := make(map[string]interface{})
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
    global := map[string]interface{}{
        "Name":          models.Config.Name,
        "Title":         "test",
        "Vers":          models.Config.Vers,
        "ServersCount":  len(servers.Servers) + len(apis.Apis),
        "Ongoing":       atks.Container.GlobalRunning(),
        "Slots":         servers.Slots()[0],
        "Users":         users.Container.Users() + site.Site.FakeUsers,
        "Remotes":       servers.Servers,
        "Session":       userSession,
        "FlashMessages": flashMessages,
    }
    for k, val := range global {
        if _, exists := data[k]; !exists {
            data[k] = val
        }
    }
    funcMap := template.FuncMap{
        "hasPerm": func(role string) bool {
            return users.HasPermission(userSession.User, role)
        },
    }

    pagePath := filepath.Join("assets/html", filepath.Join(file...))
    partials := []string{
        "assets/html/nav.html",
        "core/footer.html",
        "assets/html/construction.html",
    }

    base := filepath.Base(pagePath)
    tmpl, err := template.
        New(base).
        Funcs(funcMap).
        ParseFiles(append([]string{pagePath}, partials...)...)
    if err != nil {
        log.Println("Render: ParseFiles:", err)
        return
    }

    if err := tmpl.ExecuteTemplate(w, base, data); err != nil {
        log.Println("Render: ExecuteTemplate:", err)
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