package attackapi

import (
	"api/core/database"
    "api/core/database/atks"
	"api/core/models/floods"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"
	"net/http"
	"strings"
)
type AttackParams struct {
    Host        string
    Port        int
    Method      string
    Duration    int
    Threads     int
    PPS         int
    Concurrents int
}
func isResolvable(target string) bool {
		if net.ParseIP(target) != nil {
			return true
		}

		parsedURL, err := url.Parse(target)
		if err != nil {
			return false
		}
		host := parsedURL.Hostname()
		if host == "" {
			host = target
		}

		// Perform DNS lookup
		addrs, err := net.LookupHost(host)
		return err == nil && len(addrs) > 0
}


func ValidateTarget(target string, blacklists []string) error {
	if net.ParseIP(target) != nil {
		if inBlacklist(target, blacklists) {
			return errors.New("target is blacklisted")
		}
		return nil
	}
	parsedURL, err := url.Parse(target)
	if err != nil {
		return errors.New("invalid target provided")
	}
	host := parsedURL.Hostname()
	if host == "" {
		host = target
	}
	if !isResolvable(host) {
		return errors.New("target is not resolvable")
	}
	if inBlacklist(target, blacklists) {
		return errors.New("target is blacklisted")
	}

	return nil
}

func inBlacklist(target string, blacklists []string) bool {
	for _, b := range blacklists {
		if b == target {
			return true
		}
	}
	return false
}

func ParseAndValidate(r *http.Request,	user *database.User, blacklists []string) (*AttackParams, error) {
    // merge URL‑query + POST‑form
    if err := r.ParseForm(); err != nil {
        return nil, err
    }

    p := &AttackParams{}

    // 1) Host
    p.Host = strings.TrimSpace(r.FormValue("host"))
    if p.Host == "" {
        return nil, errors.New("host is required")
    }
    if err := ValidateTarget(p.Host, blacklists); err != nil {
        return nil, err
    }

    // 2) Port
    port, err := strconv.Atoi(r.FormValue("port"))
    if err != nil {
        return nil, errors.New("invalid port")
    }
    if port < 1 || port > 65535 {
        return nil, errors.New("port out of range")
    }
    p.Port = port
	
	durStr := r.FormValue("duration")
	if durStr == "" {
		// fallback for GET
		durStr = r.FormValue("time")
	}
    duration, err := strconv.Atoi(durStr)
    if err != nil {
        return nil, errors.New("invalid duration")
    }
    if duration <= 0 || duration > user.Duration {
        return nil, errors.New("duration exceeds maximum allowed")
    }
    p.Duration = duration

    p.Method = r.FormValue("method")
    if floods.New(p.Method) == nil {
        return nil, errors.New("invalid method")
    }

    if v := r.FormValue("threads"); v != "" {
        if t, err := strconv.Atoi(v); err == nil {
            p.Threads = t
        } else {
            return nil, errors.New("invalid threads")
        }
    }

    if v := r.FormValue("pps"); v != "" {
        if pp, err := strconv.Atoi(v); err == nil {
            p.PPS = pp
        } else {
            return nil, errors.New("invalid pps")
        }
    }

    p.Concurrents = 1
    ongoing, _ := atks.Container.GetRunning(user)
    if len(ongoing) >= user.Concurrents {
        return nil, errors.New("maximum running attacks reached")
    }
    if v := r.FormValue("concurrents"); v != "" {
        parts := strings.Split(v, ".")[0]
        c, err := strconv.Atoi(parts)
        if err != nil {
            return nil, errors.New("invalid concurrents")
        }
        if c < 1 || len(ongoing)+c > user.Concurrents {
            return nil, errors.New("you're trying to attack with more concurrents than you have available")
        }
        p.Concurrents = c
    }

    return p, nil
}

func SaveToDB(user *database.User, flood *floods.Attack, conns int) ([]int, error) {
    var ids []int
    for i := 0; i < conns; i++ {
        id, err := atks.Container.NewAttack(user, flood)
        if err != nil {
            return nil, fmt.Errorf("database error occurred: %w", err)
        }
        ids = append(ids, id)
        time.Sleep(500 * time.Microsecond)
    }
    return ids, nil
}
