package site

import (
	"database/sql"
	"encoding/base64"
	"api/core/database"
	"log"
)

func (conn *Instance) GetApiKey(username string) (string, error) {
	stmt, err := conn.conn.Prepare("SELECT `api` FROM `users` WHERE `username` = ?")
	if err != nil {
		return "", err 
	}
	defer stmt.Close()
	var apiKey string
	err = stmt.QueryRow(username).Scan(&apiKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUserNotFound 
		}
		return "", err 
	}
	decodedKey, err := base64.StdEncoding.DecodeString(apiKey)
    if err != nil {
        return "", err // Return an error if decoding fails
    }
	apiKey = string(decodedKey)
	return apiKey, nil 
}

func (conn *Instance) IsApiKey(key string, user *database.User) bool {
    stmt, err := conn.conn.Prepare("SELECT `api` FROM `users` WHERE `username` = ?")
    if err != nil {
        log.Printf("Error preparing query: %v", err)
        return false
    }
    defer stmt.Close()

    var apiKey string
    err = stmt.QueryRow(user.Username).Scan(&apiKey)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("No API key found for user %s", user.Username)
            return false
        }
        log.Printf("Error querying API key: %v", err)
        return false
    }

    // Decode the base64-encoded API key from the database
    decodedKey, err := base64.StdEncoding.DecodeString(apiKey)
    if err != nil {
        log.Printf("Error decoding API key: %v", err)
        return false
    }

    // Compare the decoded key with the given key
    return string(decodedKey) == key
}
func (conn *Instance) GetReqs(username string) (int, int, error) {
	stmt, err := conn.conn.Prepare("SELECT `apiReqs`, `apiFails` FROM `users` WHERE `username` = ?")
	if err != nil {
		return 0, 0, err
	}
	defer stmt.Close()

	var apiReqs, apiFails int
	err = stmt.QueryRow(username).Scan(&apiReqs, &apiFails)
	if err != nil {
		return 0, 0, err
	}

	return apiReqs, apiFails, nil
}

func (conn *Instance) SaveReqs(username string, success bool) error {
	// Base update: always increment apiReqs
	query := "UPDATE `users` SET `apiReqs` = `apiReqs` + 1"

	if !success {
		query += ", `apiFails` = `apiFails` + 1"
	}

	query += " WHERE `username` = ?"

	stmt, err := conn.conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username)
	return err
}