package database

import (
	"api/core/models"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

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
		db, err = sql.Open("sqlite3", "assets/database.db")
		if err != nil {
			return fmt.Errorf("could not open SQLite connection: %v", err)
		}
	}

	// Check database connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("could not ping the database: %v", err)
	}

	Container.conn = db
	log.Println("New(): successfully connected to database")
	return nil
}
