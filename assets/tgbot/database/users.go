package database

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID         int
	Username   string
	Telegram   string
	Membership string
	Balance    int
}

// GetUser by username
func (conn *Instance) GetUser(username string) (*User, error) {
	stmt, err := conn.conn.Prepare("SELECT id, username, membership, balance, telegram FROM users WHERE telegram = ?")
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %v", err)
	}
	defer stmt.Close()

	var user User
	err = stmt.QueryRow(username).Scan(&user.ID, &user.Username, &user.Membership, &user.Balance, &user.Telegram)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	return &user, nil
}

// AddUser
func (conn *Instance) AddUser(user *User) error {
	stmt, err := conn.conn.Prepare(`
		INSERT OR REPLACE INTO users (username, telegram, membership, balance)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("error preparing insert statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Telegram, user.Membership, user.Balance)
	if err != nil {
		return fmt.Errorf("error inserting or replacing user into database: %v", err)
	}

	return nil
}