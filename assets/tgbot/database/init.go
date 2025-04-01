package database

import (
	"bot/models"
	"database/sql"
	"api/core/models/log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Container = new(Instance)
)

type Query interface{ Scan(...any) error }

type Instance struct {
	Connected time.Time

	conn *sql.DB
}

func New() error {
	Container.Connected = time.Now()
	db, err := sql.Open("sqlite3", models.Config.DbPath)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	Container.conn = db
	lologgger.Println("New(): succesfully connected to database")
	return nil
}
