package panelapi

import (
	"api/core/master/sessions"
	"api/core/models/floods"
	"api/core/models/server"
    "api/core/models"
	"encoding/json"
	"net/http"
	"strings"
    "fmt"
)

func init() {
    Route.NewSub(server.NewRoute("/methods", func(w http.ResponseWriter, r *http.Request) {
        if strings.ToLower(r.Method) == "post" {
            ok, user := sessions.IsLoggedIn(w, r)
            if !ok {
                http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
                return
            }
            type method struct {
                Description string `json:"description"`
                ID          int    `json:"id"`
                Method      string `json:"method"`
                PanelMethod string `json:"panel_method"`
                Subnet      string `json:"subnet"`
                Mtype       int    `json:"mtype"`
                VIP         string `json:"vip"`
            }
            type status struct {
                Status  string    `json:"status"`
                Methods []*method `json:"methods"`
            }
            var s = &status{
                Status:  "success",
                Methods: make([]*method, 0),
            }
            if !user.HasPermission("basic") && !user.HasPermission("vip") && !user.HasPermission("admin") && !user.HasPermission("api") {
                s.Methods = append(s.Methods, &method{
                    PanelMethod: "PLEASE BUY A PLAN",
                })
            }
            if user.HasPermission("vip") || user.HasPermission("admin") {
                for name, meth := range floods.Methods {
                    s.Methods = append(s.Methods, &method{
                        Description: meth.Description,
                        Method:      name,
                        PanelMethod: meth.Name,
                        ID:          0,
                        Subnet:      models.Config.Methods[fmt.Sprintf("subnet%d", meth.Subnet)],
                        VIP:         "VIP",  
                        Mtype:       meth.Mtype,
                    })
                }
            } else if user.HasPermission("basic") {
                for name, meth := range floods.Methods {
                    if meth.VIP == false {
                        s.Methods = append(s.Methods, &method{
                            Description: meth.Description,
                            Method:      name,
                            PanelMethod: meth.Name,
                            ID:          0,
                            Subnet:      models.Config.Methods[fmt.Sprintf("subnet%d", meth.Subnet)],
                            VIP:         "BASIC",  
                            Mtype:       meth.Mtype,
                        })
                    }
                }
            }
            json.NewEncoder(w).Encode(s)
            return
        } else {
            w.Write([]byte("404 page not found"))
            w.WriteHeader(404)
        }
    }))
}
