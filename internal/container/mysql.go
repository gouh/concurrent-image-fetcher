package container

import (
	"concurrent-image-fetcher/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func NewMysql(config *config.DatabaseConfig) *sql.DB {
	// Open database connection
	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	// Set maximum number of connections in idle connection pool
	db.SetMaxIdleConns(10)

	// Set maximum number of open connections to the database
	db.SetMaxOpenConns(100)

	// Set the maximum amount of time a connection may be reused
	db.SetConnMaxLifetime(time.Hour)

	return db
}
