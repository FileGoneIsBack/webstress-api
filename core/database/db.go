package database

import (
	"database/sql"
	"time"
)

var (
	Container = new(Instance)
)

type Query interface{ Scan(...any) error }

type Instance struct {
	Connected time.Time

	conn *sql.DB
}
