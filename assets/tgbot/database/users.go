package database

import (
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	ID          int
	Username    string
	Api         []byte
	Roles       string
	Expiry      int
	Membership  string
	Concurrents int
	Servers     int
	Duration    int
	Balance     int
	ApiReqs     int
	ApiFails    int
	Tele        int64
}

// GetUser by username
func (conn *Instance) GetUser(tele int64) (*User, error) {
	log.Printf("\n\n %d\n\n", tele)

	stmt, err := conn.Conn.Prepare(`
	SELECT id, username, api, roles, expiry, membership, concurrents, servers, duration, balance, apiReqs, apiFails, tele
	FROM users
	WHERE tele = ?`)
	if err != nil {
		return nil, fmt.Errorf("\n\nerror preparing query: %v", err)
	}
	defer stmt.Close()

	var user User
	err = stmt.QueryRow(tele).Scan(
		&user.ID,
		&user.Username,
		&user.Api,
		&user.Roles,
		&user.Expiry,
		&user.Membership,
		&user.Concurrents,
		&user.Servers,
		&user.Duration,
		&user.Balance,
		&user.ApiReqs,
		&user.ApiFails,
		&user.Tele,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("\n\nerror executing query: %v", err)
	}

	return &user, nil
}
