package sessions

import (
	"api/core/database"
	"time"
)

var (
	Sessions = map[string]Session{}
)

// Session is used to store the user & expiry
type Session struct {
     *database.User
	Expiry time.Time
	Flashes []interface{}
}

// IsExpired is used to determine if the Session has expired
func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

func Count() int {
    return len(Sessions)
}
