package database

import (
	"api/core/models/plans"
	"api/core/models/ranks"
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

type Role string
type User struct {
	ID                                      int
	Username                                string
	Key, Salt                               []byte
	ranks, Membership                       string
	Ranks                                   []*ranks.Rank
	Concurrents, Servers, Duration, Balance int
	Expiry                                  int64
}

var (
	ErrDuplicateUser = errors.New("duplicate user")
	ErrUserNotFound  = errors.New("user couldn't be found in the database")
	ErrInvalidInput  = errors.New("invalid ticket data")
)

func (conn *Instance) NewUser(user *User) (err error) {
	if user, err := conn.GetUser(user.Username); err == nil && user != nil {
		return ErrDuplicateUser
	}
	user.Salt = NewSalt(16)
	user.Key = NewHash(user.Key, user.Salt)
	user.ranks = user.NewRoles()
	stmt, err := conn.conn.Prepare("INSERT INTO `users` (`id`, `username`, `key`, `salt`, `roles`, `expiry`, `concurrents`, `servers`, `duration`, `balance`, `membership`) VALUES (NULL, ?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(user.Username, user.Key, user.Salt, user.ranks, user.Expiry, user.Concurrents, user.Servers, user.Duration, user.Balance, user.Membership); err != nil {
		return err
	}
	return
}

func (conn *Instance) UserUpdateAddon(username string, balance, duration, concurrents int, roles []*ranks.Rank, addon *plans.Addon) error {
	logger.Printf("Updating addon for user %s\n", username)

	// Convert ranks to JSON string if addon is of type rank
	var encodedJSON string
	if addon.Type == "rank" {
		ranksJSON, err := json.Marshal(roles)
		if err != nil {
			logger.Println("Error marshalling ranks to JSON:", err)
			return err
		}
		encodedJSON = base64.RawStdEncoding.EncodeToString(ranksJSON)
	}

	// Prepare the SQL statement
	query := "UPDATE `users` SET `balance` = ?, `duration` = ?, `concurrents` = ?, `roles` = ? WHERE `username` = ?"
	stmt, err := conn.conn.Prepare(query)
	if err != nil {
		logger.Println("Error preparing SQL statement:", err)
		return err
	}
	defer stmt.Close()

	// Set values based on the addon type
	var newDuration, newConcurrents int
	var newRoles string

	if addon.Type == "time" {
		newDuration = duration + addon.Value
		newConcurrents = concurrents // Keep concurrents unchanged
		newRoles = ""                // No roles update
	} else if addon.Type == "concurrents" {
		newDuration = duration // Keep duration unchanged
		newConcurrents = concurrents + addon.Value
		newRoles = "" // No roles update
	} else if addon.Type == "rank" {
		newDuration = duration       // Keep duration unchanged
		newConcurrents = concurrents // Keep concurrents unchanged
		newRoles = encodedJSON       // Update roles
	}

	// Execute the update query with the calculated values
	_, err = stmt.Exec(balance-addon.Price, newDuration, newConcurrents, newRoles, username)
	if err != nil {
		logger.Println("Error executing SQL statement:", err)
		return err
	}

	return nil
}

func (conn *Instance) GetUser(user string) (*User, error) {
	stmt, err := conn.conn.Prepare("SELECT `id`, `username`, `key`, `salt`, `roles`, `expiry`, `concurrents`, `servers`, `duration`, `balance`, `membership` FROM `users` where `username` = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return conn.scanUser(stmt.QueryRow(user))
}

func (conn *Instance) GetUserID(username string) (int, error) {
	var userID int
	stmt, err := conn.conn.Prepare("SELECT `id` FROM `users` WHERE `username` = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// Execute the query and scan the result into userID
	err = stmt.QueryRow(username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrUserNotFound
		}
		return 0, err
	}

	return userID, nil
}

func (conn *Instance) GetUserByID(id int) (*User, error) {
	stmt, err := conn.conn.Prepare("SELECT `id`, `username`, `key`, `salt`, `roles`, `expiry`, `concurrents`, `servers`, `duration`, `balance`, `membership` FROM `users` where `id` = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return conn.scanUser(stmt.QueryRow(id))
}

func (conn *Instance) GetUsers() ([]*User, error) {
	stmt, err := conn.conn.Prepare("SELECT `id`, `username`, `key`, `salt`, `roles`, `expiry`, `concurrents`, `servers`, `duration`, `balance`, `membership` FROM `users`")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	var users []*User
	for rows.Next() {
		user, err := conn.scanUser(rows)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func (conn *Instance) UpdateUserPlan(user *User, plan *plans.Plan) error {
	stmt, err := conn.conn.Prepare("UPDATE `users` SET `roles` = ?, `expiry` = ?, `concurrents` = ?, `servers` = ?, `duration` = ?, `balance` = ?, `membership` = ? WHERE `username` = ?")
	if err != nil {
		return err
	}

	// Check if the user already has the "admin" role before clearing it
	hasAdminRole := false
	for _, rank := range user.Ranks {
		if rank.Name == "admin" {
			hasAdminRole = true
			break
		}
	}

	user.Ranks = []*ranks.Rank{}
	if hasAdminRole {
		user.Ranks = append(user.Ranks, ranks.GetRole("admin", true))
	}
	if user.HasPermission("admin") {
		user.Ranks = append(user.Ranks, ranks.GetRole("api", true), ranks.GetRole("vip", true), ranks.GetRole("api", true), ranks.GetRole("cnc", true), ranks.GetRole("admin", true))
	} else if user.HasPermission("basic") {
		user.Ranks = append(user.Ranks, ranks.GetRole("basic", false))
	} else if plan.API && plan.VIP {
		user.Ranks = append(user.Ranks, ranks.GetRole("api", true), ranks.GetRole("vip", true))
		user.Membership = "premium"
	} else if plan.API {
		user.Ranks = append(user.Ranks, ranks.GetRole("api", true), ranks.GetRole("basic", true))
		user.Membership = "API"
	} else if plan.VIP {
		user.Ranks = append(user.Ranks, ranks.GetRole("vip", true))
		user.Membership = "VIP"
	} else if !plan.VIP {
		user.Ranks = append(user.Ranks, ranks.GetRole("basic", true))
		user.Membership = "basic"
	}

	user.ranks = user.NewRoles() 
	user.Balance -= plan.Price  
	if _, err := stmt.Exec(user.ranks, time.Now().Add((time.Duration(plan.Expiry)*time.Hour)*24).Unix(), plan.Conns, 5, plan.Duration, user.Balance, user.Membership, user.Username); err != nil {
		return err
	}

	return nil
}

func (user *User) GetKey() []byte {
	return user.Key
}
func (conn *Instance) scanUser(query Query) (*User, error) {
	user := new(User)
	if err := query.Scan(
		&user.ID,
		&user.Username,
		&user.Key, &user.Salt, &user.ranks, &user.Expiry, &user.Concurrents, &user.Servers, &user.Duration, &user.Balance, &user.Membership,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	if err := user.Sync(); err != nil {
		return user, err
	}
	return user, nil
}

func (user *User) IsKey(key []byte) bool {
	return bytes.Equal(NewHash(key, user.Salt), user.Key) && true
}

func (conn *Instance) Users() (users int) {
	stmt, err := conn.conn.Prepare("SELECT * from `users`")
	if err != nil {
		return 0
	}
	defer stmt.Close()
	result, err := stmt.Query()
	if err != nil {
		logger.Println("GlobalUsers(): error occured while executing statement \"" + err.Error() + "\"")
		return 0
	}
	for result.Next() {
		users++
	}
	return
}

func (conn *Instance) UpdateUser(user *User) error {
	stmt, err := conn.conn.Prepare("UPDATE `users` SET `roles` = ?, `expiry` = ?, `concurrents` = ?, `servers` = ?, `duration` = ?, `balance` = ? WHERE `username` = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	user.ranks = user.NewRoles()
	// Execute the update query
	if _, err := stmt.Exec(user.ranks, user.Expiry, user.Concurrents, user.Servers, user.Duration, user.Balance, user.Username); err != nil {
		return err
	}

	return nil
}

func (conn *Instance) DeleteUser(username string, userID int) error {
	// Prepare the delete statement
	stmt, err := conn.conn.Prepare("DELETE FROM `users` WHERE `username` = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the delete statement
	res, err := stmt.Exec(username)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	// Prepare the delete statement
	stmt1, err := conn.conn.Prepare("DELETE FROM `tickets` WHERE `user_id` = ?")
	if err != nil {
		return err
	}
	defer stmt1.Close()

	// Execute the delete statement
	res, err = stmt1.Exec(userID)
	if err != nil {
		return err
	}

	return nil
}

func (conn *Instance) UserData(row *sql.Row) (*User, error) {
	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Key, &u.Salt, &u.Ranks, &u.Expiry, &u.Concurrents, &u.Servers, &u.Duration, &u.Balance, &u.Membership)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

