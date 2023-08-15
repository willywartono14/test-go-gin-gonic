package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/willywartono14/test-go-gin-gonic/config"
)

func Init() *sql.DB {
	var dbConfig = config.Get().Database

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.DbName, dbConfig.SslMode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("[DATABASE] Failed to initialize database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("[DATABASE] Failed to ping database: %v", err)
		if err = db.Close(); err != nil {
			log.Fatalf("[DATABASE] Failed to close database: %v", err)
		}
	}

	return db
}
