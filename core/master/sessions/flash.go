package sessions

import (
	"errors"
	"net/http"
)

type FlashMessage struct {
    //time    time.Time
	From    string
    Message string
}

// SetFlash sets a flash message into the current session, limiting the number to 4.
func SetFlash(w http.ResponseWriter, r *http.Request, message, from string) error {
    cookie, err := r.Cookie("session-token")
    if err != nil {
        return err
    }
    sessionToken := cookie.Value
    session, exists := Sessions[sessionToken]
    if !exists {
        return errors.New("session not found")
    }

    if len(session.Flashes) >= 4 {
        session.Flashes = session.Flashes[1:]
    }

    session.Flashes = append(session.Flashes, FlashMessage{
        Message: message,
        From: from, 
        //Time: 
    })

    Sessions[sessionToken] = session
    return nil
}

func GetFlash(w http.ResponseWriter, r *http.Request) []FlashMessage {
    cookie, err := r.Cookie("session-token")
    if err != nil {
        return nil
    }
    sessionToken := cookie.Value
    session, exists := Sessions[sessionToken]
    if !exists {
        return nil
    }

    var messages []FlashMessage
    for _, flash := range session.Flashes {
        if msg, ok := flash.(FlashMessage); ok {
            messages = append(messages, msg)
        }
    }
    return messages
}

// ClearFlashes clears the flash messages after they're rendered.
func ClearFlashes(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session-token")
	if err != nil {
		return err
	}
	sessionToken := cookie.Value
	session, exists := Sessions[sessionToken]
	if !exists {
		return errors.New("session not found")
	}

	session.Flashes = nil
	Sessions[sessionToken] = session
	return nil
}

func BroadcastFlash(message interface{}) {
	for token, session := range Sessions {
		if len(session.Flashes) >= 4 {
			session.Flashes = session.Flashes[1:]
		}
		session.Flashes = append(session.Flashes, message)
		Sessions[token] = session
	}
}