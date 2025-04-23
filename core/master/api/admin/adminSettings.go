package adminapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"api/core/database/site"
	"api/core/database/users"
	"api/core/master/sessions"
	"api/core/models/server"
)

// helper to write a JSON error
func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func init() {
	Route.NewSub(server.NewRoute("/settingsEnd", func(w http.ResponseWriter, r *http.Request) {
		ok, session := sessions.IsLoggedIn(w, r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		if !users.HasPermission(session.User, "admin") {
			http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
			return
		}
		
		switch r.Method {
		case http.MethodGet:
			getSettings(w)
		case http.MethodPost:
			updateSettings(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		}
	}))
}

func getSettings(w http.ResponseWriter) {
	if site.Site == nil {
		jsonError(w, "settings not loaded", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(site.Site)
}

func updateSettings(w http.ResponseWriter, r *http.Request) {
  var updates map[string]interface{}
  if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
    jsonError(w, "invalid JSON payload", http.StatusBadRequest)
    return
  }

  for field, raw := range updates {
    lf := strings.ToLower(field)
    if !site.AllowedFields[lf] {
      jsonError(w, "unknown field "+lf, http.StatusBadRequest)
      return
    }

    var val interface{}
    switch lf {
    case "freeuser", "cnc", "api":
      b, ok := raw.(bool)
      if !ok {
        jsonError(w, lf+" must be true or false", http.StatusBadRequest)
        return
      }
      val = b

    case "freeuser_cons", "freeuser_time", "fake_users", "fake_ongoing":
      if f, ok := raw.(float64); ok {
        val = int(f)
      } else {
        jsonError(w, lf+" must be a number", http.StatusBadRequest)
        return
      }

    // ← handle roles *before* default:
    case "freeuser_roles":
      arr, ok := raw.([]interface{})
      if !ok {
        jsonError(w, "freeuser_roles must be an array of strings", http.StatusBadRequest)
        return
      }
      roles := make([]string, len(arr))
      for i, v := range arr {
        s, ok := v.(string)
        if !ok {
          jsonError(w, "freeuser_roles must be an array of strings", http.StatusBadRequest)
          return
        }
        roles[i] = s
      }
      // join into a CSV for storage
      val = strings.Join(roles, ",")

    default:
      // everything else is stored as string
      s, ok := raw.(string)
      if !ok {
        jsonError(w, lf+" must be a string", http.StatusBadRequest)
        return
      }
      val = s
    }

    if err := site.Container.UpdateSiteSettingField(lf, val); err != nil {
      jsonError(w, "failed to update "+lf+": "+err.Error(), http.StatusInternalServerError)
      return
    }
  }

  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(map[string]string{"status":"ok"})
}