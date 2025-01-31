package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
	"strings"
)

type Invite struct {
	Token 				int		
	Expiry              time.Time
	Uses   				[]string
}

func (conn *Instance) LogUserUsage(token int, username string) error {
    // Fetch the current usernames for this token
    var currentUsernames sql.NullString
    stmt, err := conn.conn.Prepare("SELECT usernames FROM tokens WHERE token = ?")
    if err != nil {
        return fmt.Errorf("error preparing the SQL query: %v", err)
    }
    defer stmt.Close()

    err = stmt.QueryRow(token).Scan(&currentUsernames)
    if err != nil && err != sql.ErrNoRows {
        return fmt.Errorf("error fetching current usernames: %v", err)
    }

    // Check if currentUsernames is valid and not NULL
    if currentUsernames.Valid {
        // Append the new username if the currentUsernames is not empty
        currentUsernames.String = currentUsernames.String + "," + username
    } else {
        currentUsernames.String = username
    }

    // Update the usernames column with the new list
    stmt, err = conn.conn.Prepare("UPDATE tokens SET usernames = ? WHERE token = ?")
    if err != nil {
        return fmt.Errorf("error preparing the SQL query: %v", err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(currentUsernames.String, token)
    if err != nil {
        return fmt.Errorf("error executing the SQL query: %v", err)
    }

    return nil
}

func (conn *Instance) GetInvite(token string, username string) (string, time.Time, error) {
    var dbToken string
    var exp time.Time

    // Prepare statement to fetch token details
    stmt, err := conn.conn.Prepare("SELECT token, exp FROM tokens WHERE token = ?")
    if err != nil {
        return "", time.Time{}, err
    }
    defer stmt.Close()

    // Execute query to fetch the token and its expiration
    err = stmt.QueryRow(token).Scan(&dbToken, &exp)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", time.Time{}, fmt.Errorf("token not found")
        }
        return "", time.Time{}, err
    }

    num, _ := strconv.Atoi(dbToken)
    if err := conn.LogUserUsage(num, username); err != nil {
        return "", time.Time{}, fmt.Errorf("error logging user usage: %v", err)
    }

    return dbToken, exp, nil
}


func (conn *Instance) AddInvite(token int, expiry time.Time) error {
	stmt, err := conn.conn.Prepare("INSERT INTO tokens (token, exp) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing the SQL query: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(token, expiry)
	if err != nil {
		return fmt.Errorf("error executing the SQL query: %v", err)
	}

	return nil
}

func (conn *Instance) GetAllInvites() ([]Invite, error) {
    stmt, err := conn.conn.Prepare("SELECT token, usernames, exp FROM tokens")
    if err != nil {
        return nil, fmt.Errorf("error preparing the SQL query: %v", err)
    }
    defer stmt.Close()

    rows, err := stmt.Query()
    if err != nil {
        return nil, fmt.Errorf("error querying the database: %v", err)
    }
    defer rows.Close()

    var invites []Invite

    // Iterate over each token entry
    for rows.Next() {
        var token int
        var usernames sql.NullString
        var exp time.Time

        // Scan token, usernames, and expiry date
        if err := rows.Scan(&token, &usernames, &exp); err != nil {
            return nil, fmt.Errorf("error scanning row: %v", err)
        }

        // Split the usernames string into a slice (handle empty usernames gracefully)
		var userList []string
		if usernames.Valid && usernames.String != "" {
			// If valid, split the comma-separated string into a slice
			userList = strings.Split(usernames.String, ",")
		} else {
			// If NULL or empty, return an empty list
			userList = []string{}
		}

        if err != nil {
            return nil, fmt.Errorf("error parsing expiry date: %v", err)
        }

        // Append the invite with users
        invites = append(invites, Invite{
            Token:  token,
            Expiry: exp,
            Uses:   userList,  // Set the Uses field to the list of usernames
        })
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating over rows: %v", err)
    }
    return invites, nil
}

func (conn *Instance) DeleteInvite(token int) error {
    stmt, err := conn.conn.Prepare("DELETE FROM tokens WHERE token = ?")
    if err != nil {
        return fmt.Errorf("error preparing the SQL query: %v", err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(token)
    if err != nil {
        return fmt.Errorf("error executing the SQL query: %v", err)
    }

    return nil
}
