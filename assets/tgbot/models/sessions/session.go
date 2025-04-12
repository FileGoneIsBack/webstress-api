package sessions

import (
    "errors"
    "sync"
	"bot/database"
)

type Session struct {
    ID        int
    Username  string
    Telegram  int64
    Balance   int
    Membership string
}

var sessions = make(map[int]*database.User)  
var sessionStore = make(map[int]*Session)
var mu sync.RWMutex

func GetSession(userID int) (*Session, error) {
    mu.RLock() 
    defer mu.RUnlock()

    session, exists := sessionStore[userID]
    if !exists {
        return nil, errors.New("session not found")
    }
    return session, nil
}

// stores a session 
func SetSession(userID int, session *Session) {
    mu.Lock() 
    defer mu.Unlock()

    sessionStore[userID] = session
}
