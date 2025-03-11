package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"api/core/models"
)

func New() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", 
		models.Config.Database.Username,   
		models.Config.Database.Password,    
		models.Config.Database.Host,        
		models.Config.Database.Database,   
	)
	Container.Connected = time.Now()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("could not open connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("could not ping the database: %v", err)
	}

	Container.conn = db
	log.Println("New(): successfully connected to MySQL database")
	return nil
}
