package database

import (
	"bot/models"
	"log"
	"time"
)

type Invite struct {
	Token  int
	Expiry time.Time
	User   string
}

func (conn *Instance) CreateInvite(user int64) (string, error) {
	token := models.GenCode()
	expiry := time.Now().Add(5 * time.Minute)

	_, err := conn.Conn.Exec(`INSERT INTO tokens (token, expiry, usernames) VALUES (?, ?, ?)`, token, expiry, user)
	if err != nil {
		return "", err
	}
	Container.ScheduleInviteDeletion(token)
	return token, nil
}

func (conn *Instance) DeleteExpiredInvites() error {
	_, err := conn.Conn.Exec(`DELETE FROM tokens WHERE expiry <= ?`, time.Now())
	return err
}

func (conn *Instance) ScheduleInviteDeletion(token string) {
	go func() {
		for i := 0; i < 10; i++ {
			log.Printf("Invite %s expires in %d seconds...", token, (9-i)*30)
			time.Sleep(30 * time.Second)
		}

		_, err := conn.Conn.Exec(`DELETE FROM tokens WHERE token = ?`, token)
		if err != nil {
			log.Printf("Failed to delete invite %s: %v", token, err)
		} else {
			log.Printf("Invite %s has been deleted.", token)
		}
	}()
}
