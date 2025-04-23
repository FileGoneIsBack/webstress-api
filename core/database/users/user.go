package users

import (
	"api/core/database"
	"api/core/models"
	"api/core/models/log"
	"api/core/models/plans"
	"api/core/models/ranks"
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
)

func (conn *Instance) NewUser(user *database.User) error {
	existingUser, err := conn.GetUser(user.Username)
	if err == nil && existingUser != nil {
		return ErrDuplicateUser
	}

	user.Api = base64.StdEncoding.EncodeToString([]byte(GenerateUserKey(models.Config.Name)))
	user.Salt = database.NewSalt(16)
	user.Key = database.NewHash(user.Key, user.Salt)
	user.RanksData = NewRoles(user)

	stmt, err := conn.conn.Prepare(`
	INSERT INTO users 
	(id, username, ` + "`key`" + `, salt, roles, expiry, concurrents, servers, duration, balance, membership, api, tele, apiReqs, apiFails) 
	VALUES 
	(NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		user.Username,
		user.Key,
		user.Salt,
		user.RanksData,
		user.Expiry,
		user.Concurrents,
		user.Servers,
		user.Duration,
		user.Balance,
		user.Membership,
		user.Api,
		user.Tele,
		user.ApiReqs,
		user.ApiFails,
	); err != nil {
		return err
	}

	_, err = conn.conn.Exec(`DELETE FROM tokens WHERE usernames = ?`, user.Tele)
	if err != nil {
		return fmt.Errorf("user created, but failed to clean up invite: %w", err)
	}
	return nil
}

func (conn *Instance) UserUpdateAddon(username string, balance, duration, concurrents int, roles []*ranks.Rank, addon *plans.Addon) error {
	log.Printf("Updating addon for user %s\n", username)

	// Convert ranks to JSON string if addon is of type rank
	var encodedJSON string
	if addon.Type == "rank" {
		ranksJSON, err := json.Marshal(roles)
		if err != nil {
			log.Println("Error marshalling ranks to JSON:", err)
			return err
		}
		encodedJSON = base64.RawStdEncoding.EncodeToString(ranksJSON)
	}

	// Prepare the SQL statement
	query := "UPDATE `users` SET `balance` = ?, `duration` = ?, `concurrents` = ?, `roles` = ? WHERE `username` = ?"
	stmt, err := conn.conn.Prepare(query)
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
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
		log.Println("Error executing SQL statement:", err)
		return err
	}

	return nil
}

func (conn *Instance) GetUser(user string) (*database.User, error) {
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

func (conn *Instance) GetUserByID(id int) (*database.User, error) {
	stmt, err := conn.conn.Prepare("SELECT `id`, `username`, `key`, `salt`, `roles`, `expiry`, `concurrents`, `servers`, `duration`, `balance`, `membership` FROM `users` where `id` = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return conn.scanUser(stmt.QueryRow(id))
}

func (conn *Instance) GetUsers() ([]*database.User, error) {
	stmt, err := conn.conn.Prepare("SELECT `id`, `username`, `key`, `salt`, `roles`, `expiry`, `concurrents`, `servers`, `duration`, `balance`, `membership` FROM `users`")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	var users []*database.User
	for rows.Next() {
		user, err := conn.scanUser(rows)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func (conn *Instance) UpdateUserPlan(user *database.User, plan *plans.Plan) error {
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
	if HasPermission(user, "admin") {
		user.Ranks = append(user.Ranks, ranks.GetRole("api", true), ranks.GetRole("vip", true), ranks.GetRole("api", true), ranks.GetRole("cnc", true), ranks.GetRole("admin", true))
	} else if HasPermission(user, "basic") {
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

	user.RanksData = NewRoles(user)
	user.Balance -= plan.Price
	if _, err := stmt.Exec(user.RanksData, time.Now().Add((time.Duration(plan.Expiry)*time.Hour)*24).Unix(), plan.Conns, 5, plan.Duration, user.Balance, user.Membership, user.Username); err != nil {
		return err
	}

	return nil
}

func GetKey(user *database.User) []byte {
	return user.Key
}

func (conn *Instance) scanUser(query Query) (*database.User, error) {
	user := new(database.User)
	if err := query.Scan(
		&user.ID,
		&user.Username,
		&user.Key, &user.Salt, &user.RanksData, &user.Expiry, &user.Concurrents, &user.Servers, &user.Duration, &user.Balance, &user.Membership,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	if err := Sync(user); err != nil {
		return user, err
	}
	return user, nil
}

func IsKey(u *database.User, raw []byte) bool {
    return bytes.Equal(database.NewHash(raw, u.Salt), u.Key)
}

func (conn *Instance) Users() (users int) {
	stmt, err := conn.conn.Prepare("SELECT * from `users`")
	if err != nil {
		return 0
	}
	defer stmt.Close()
	result, err := stmt.Query()
	if err != nil {
		log.Println("GlobalUsers(): error occured while executing statement \"" + err.Error() + "\"")
		return 0
	}
	for result.Next() {
		users++
	}
	return
}

func (conn *Instance) UpdateUser(user *database.User) error {
	stmt, err := conn.conn.Prepare("UPDATE `users` SET `roles` = ?, `concurrents` = ?, `servers` = ?, `duration` = ?, `balance` = ? WHERE `username` = ?")
	if err != nil {
		return err
	} 
	defer stmt.Close()
	user.RanksData = NewRoles(user)
	// Execute the update query without updating expiry
	if _, err := stmt.Exec(user.RanksData, user.Concurrents, user.Servers, user.Duration, user.Balance, user.Username); err != nil {
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

func (conn *Instance) UserData(row *sql.Row) (*database.User, error) {
	var u database.User
	err := row.Scan(&u.ID, &u.Username, &u.Key, &u.Salt, &u.Ranks, &u.Expiry, &u.Concurrents, &u.Servers, &u.Duration, &u.Balance, &u.Membership)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (conn *Instance) UpdateUserKey(username string, newKey []byte) error {
	encodedKey := base64.StdEncoding.EncodeToString(newKey)
	stmt, err := conn.conn.Prepare("UPDATE `users` SET `api` = ? WHERE `username` = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(encodedKey, username)
	if err != nil {
		return err
	}

	return nil
}

func GenerateUserKey(siteName string) string {
	siteName = strings.ReplaceAll(siteName, " ", "-")
	randomString := generateRandomString(6)
	return fmt.Sprintf("%s-%s", siteName, randomString)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result strings.Builder
	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result.WriteByte(charset[randomIndex.Int64()])
	}
	return result.String()
}
