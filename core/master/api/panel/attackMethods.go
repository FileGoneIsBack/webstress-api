package panelapi

import (
    "api/core/database/users"
	"api/core/master/sessions"
	"api/core/models/floods"
	"api/core/models/server"
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
            if !users.HasPermission(user.User, "basic") && !users.HasPermission(user.User, "vip") && !users.HasPermission(user.User, "admin") && !users.HasPermission(user.User, "api") {
                s.Methods = append(s.Methods, &method{
                    PanelMethod: "PLEASE BUY A PLAN",
                })
            }
            if users.HasPermission(user.User, "vip") || users.HasPermission(user.User, "admin") {
                for name, meth := range floods.Methods {
                    s.Methods = append(s.Methods, &method{
                        Description: meth.Description,
                        Method:      name,
                        PanelMethod: meth.Name,
                        ID:          0,
                        Subnet:      fmt.Sprintf("subnet%d", meth.Subnet),
                        VIP:         "VIP",  
                        //Mtype:       meth.Mtype,
                    })
                }
            } else if users.HasPermission(user.User, "basic") {
                for name, meth := range floods.Methods {
                    if !meth.VIP {
                        s.Methods = append(s.Methods, &method{
                            Description: meth.Description,
                            Method:      name,
                            PanelMethod: meth.Name,
                            ID:          0,
                            Subnet:      fmt.Sprintf("subnet%d", meth.Subnet),
                            VIP:         "BASIC",  
                            //Mtype:       meth.Mtype,
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
