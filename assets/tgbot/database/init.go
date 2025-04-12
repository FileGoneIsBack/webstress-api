package database

import (
	"bot/models"
	"database/sql"
	"time"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
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

	var db *sql.DB
	var err error

	if models.Config.Secure {
		// Use MySQL
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
			models.Config.Database.Username,
			models.Config.Database.Password,
			models.Config.Database.Host,
			models.Config.Database.Database,
		)
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return fmt.Errorf("could not open MySQL connection: %v", err)
		}
	} else {
		// Use SQLite
		db, err = sql.Open("sqlite3", models.Config.DbPath)
		if err != nil {
			return fmt.Errorf("could not open SQLite connection: %v", err)
		}
	}

	// Check database connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("could not ping the database: %v", err)
	}

	Container.conn = db
	fmt.Println("New(): successfully connected to database")
	return nil
}
