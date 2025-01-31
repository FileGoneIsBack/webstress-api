package sessions

import (
	"api/core/database"
	websessions "api/core/master/sessions"
	"api/core/net/term"
	"net"
	"net/http"
	"sync"
	"time"
	"log"
	"os"
	"github.com/google/uuid"
)

var Sessions = make(map[int64]*Session)
var SessionMutex sync.Mutex
var logger  = log.New(os.Stderr, "[main] ", log.Ltime|log.Lshortfile)

type Session struct {
	ID     int64
	User   *database.User
	cookie *http.Cookie
	Conn   net.Conn
	Term   *term.Term
	Chat   bool
}

func IsLoggedIn(userID int64) (*Session, bool) {
	SessionMutex.Lock()
	defer SessionMutex.Unlock()

	logger.Println("Checking session for user ID:", userID)
	session, exists := Sessions[userID]
	if exists {
		logger.Println("Session found for user ID:", userID)
		return session, true
	} else {
		logger.Println("Session not found for user ID:", userID)
		return nil, false
	}
}

func Count() int {
	return len(Sessions)
}

func (s *Session) CreateCookie() {
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(30 * time.Minute)
	websessions.Sessions[sessionToken] = websessions.Session{
		User:   s.User,
		Expiry: expiresAt,
	}
	s.cookie = &http.Cookie{
		Name:    "session-token",
		Value:   sessionToken,
		Expires: expiresAt,
	}
}

func (s *Session) Cookie() *http.Cookie {
	return s.cookie
}

func RemoveSession(id int64) bool {
	SessionMutex.Lock()
	delete(Sessions, id)
	SessionMutex.Unlock()
	return true
}
