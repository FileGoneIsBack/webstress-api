package paymentsapi

import (
    "api/core/database/users"
    "api/core/master/sessions"
    "api/core/models/server"
    "api/core/models/ranks"
    "api/core/models/plans"
    "encoding/json"
    "net/http"
    "strings"
    "fmt"
)

func init() {
    Route.NewSub(server.NewRoute("/buyaddon", func(w http.ResponseWriter, r *http.Request) {
        type Status struct {
            Status  string `json:"status"`
            Message string `json:"message"`
        }

        switch strings.ToLower(r.Method) {
        case "post":
            ok, user := sessions.IsLoggedIn(w, r)
            if !ok {
                http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
                return
            }

            addonVal := r.PostFormValue("addon_name")

            var addon *plans.Addon
            var addonRanks []*ranks.Rank

            if addonVal == "time" {
                addon = plans.Addons["time"]
                if user.Balance >= addon.Price {
                    user.Duration += addon.Value
                    if err := users.Container.UserUpdateAddon(user.Username, user.Balance, user.Duration, user.Concurrents, addonRanks, addon); err != nil {
                        json.NewEncoder(w).Encode(&Status{Status: "error", Message: "Failed to add time."})
                        return
                    }
                    json.NewEncoder(w).Encode(&Status{Status: "success", Message: "Successfully added time."})
                    sessions.SetFlash(w, r,  fmt.Sprintf("succesfully purchased %s", addonVal), "")
                } else {
                    json.NewEncoder(w).Encode(&Status{Status: "error", Message: "Insufficient Balance!"})
                }
            
            } else if addonVal == "concurrents" {
                addon = plans.Addons["concurrents"]
                if user.Balance >= addon.Price {
                    user.Concurrents += addon.Value
                    if err := users.Container.UserUpdateAddon(user.Username, user.Balance, user.Duration, user.Concurrents, addonRanks, addon); err != nil {
                        json.NewEncoder(w).Encode(&Status{Status: "error", Message: "Failed to add concurrent connections."})
                        return
                    }
                    sessions.SetFlash(w, r,  fmt.Sprintf("succesfully purchased %s", addonVal), "")
                    json.NewEncoder(w).Encode(&Status{Status: "success", Message: "Successfully added concurrent connections."})
                } else {
                    json.NewEncoder(w).Encode(&Status{Status: "error", Message: "Insufficient Balance!"})
                }
            
            } else {
                // Handling other adons (roles)
                addonRanks = []*ranks.Rank{
                    ranks.GetRole(addonVal, true), // Get the role based on addon name
                }
            
                addon = plans.Addons[addonVal] // Assuming you have other roles in the `Addons` map
            
                if user.Balance >= addon.Price {
                    if err := users.Container.UserUpdateAddon(user.Username, user.Balance, user.Duration, user.Concurrents, addonRanks, addon); err != nil {
                        json.NewEncoder(w).Encode(&Status{Status: "error", Message: "Error updating rank."})
                        return
                    }
                    sessions.SetFlash(w, r,  fmt.Sprintf("succesfully purchased %s", addonVal), "")
                    json.NewEncoder(w).Encode(&Status{Status: "success", Message: "Successfully purchased addon " + addonVal})
                } else {
                    json.NewEncoder(w).Encode(&Status{Status: "error", Message: "Insufficient Balance!"})
                }
            }
            

            return
        }
    }))
}