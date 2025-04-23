package database

import (
	"database/sql"
	"time"
	"api/core/models/ranks"
)

var (
	Container = new(Instance)
)

type Query interface{ Scan(...any) error }

type Instance struct {
	Connected time.Time
	Conn *sql.DB
}

type User struct {
	ID                                      int
	Username, Api                           string
	Key, Salt                               []byte
	RanksData                               string
	Membership                              string
	Ranks                                   []*ranks.Rank
	Concurrents, Servers, Duration, Balance int
	Expiry, Tele                            int64
	Flashes                                 []FlashMessage
	ApiReqs, ApiFails                       int
}

type FlashMessage struct {
	Message string `json:"message"`
	Time    int64  `json:"time"`
}

type Role string


func (db *Instance) Prepare(query string) (*sql.Stmt, error) {
    return db.Conn.Prepare(query) 
}
func (db *Instance) Exec(query string, args ...any) (sql.Result, error) {
    return db.Conn.Exec(query, args...)
}
func (db *Instance) QueryRow(query string, args ...any) *sql.Row {
    return db.Conn.QueryRow(query, args...)
}
func (db *Instance) Query(query string, args ...any) (*sql.Rows, error) {
    return db.Conn.Query(query, args...)
}