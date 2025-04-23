package users

import (
	"fmt"
	"time"
)

type Invite struct {
	Token 				int		
	Expiry              time.Time
	Uses   				[]string
}

func (conn *Instance) CheckInvite(token string) (inviteUser int64, err error) {
	var expiryStr string
	err = conn.conn.QueryRow(`SELECT expiry, usernames FROM tokens WHERE token = ?`, token).Scan(&expiryStr, &inviteUser)
	if err != nil {
		return 0, fmt.Errorf("token not found or invalid: %w", err)
	}
	expiry, err := time.Parse("2006-01-02 15:04:05", expiryStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse expiry time: %w", err)
	}
	if time.Now().After(expiry) {
		return 0, fmt.Errorf("token expired")
	}
    fmt.Println("Invite user:", inviteUser)
	return inviteUser, nil
}

