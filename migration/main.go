package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/willywartono14/test-go-gin-gonic/config"
	"github.com/willywartono14/test-go-gin-gonic/database"
)

const (
	dialect string = "postgres"
)

func main() {
	err := config.Init("config.yaml")
	if err != nil {
		log.Fatalf("[CONFIG] Failed to initialize config: %v", err)
	}

	db := database.Init()
	migration := &dbMigration{
		db:      db,
		dialect: dialect,
		source: &migrate.FileMigrationSource{
			Dir: "migration/files",
		},
	}

	handleCommand(migration)
}

type dbMigration struct {
	db      *sql.DB
	dialect string
	source  *migrate.FileMigrationSource
}

func (m *dbMigration) Up() {
	n, err := migrate.Exec(m.db, m.dialect, m.source, migrate.Up)
	if err != nil {
		log.Fatalf("[MIGRATION] Failed to execute migration up: %v", err)
	}

	log.Printf("[MIGRATION] Applying %d migration(s)", n)
}

func (m *dbMigration) Down() {
	n, err := migrate.Exec(m.db, m.dialect, m.source, migrate.Down)
	if err != nil {
		log.Fatalf("[MIGRATION] Failed to execute migration down: %v", err)
	}

	log.Printf("[MIGRATION] Rolling back %d migration(s)", n)
}

func (m *dbMigration) New(fileName string) {
	content := `-- +migrate Up

-- +migrate Down

`

	fileName = fmt.Sprintf("%d_%s.sql", time.Now().UnixMilli(), fileName)
	migrationFilePath := fmt.Sprintf("migration/files/%s", fileName)
	f, err := os.Create(migrationFilePath)
	if err != nil {
		log.Fatalf("[MIGRATION] Failed to create migration file: %v", err)
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		log.Fatalf("[MIGRATION] Failed to write migration file: %v", err)
	}
}

func handleCommand(migration *dbMigration) {
	const (
		migrateUpCommand   = "migrate-up"
		migrateDownCommand = "migrate-down"
		migrateNewCommand  = "migrate-new"
	)
	permittedCommands := []string{migrateUpCommand, migrateDownCommand, migrateNewCommand}
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("[MIGRATION] Require a command. Possible command %v", permittedCommands)
	}

	command := args[1]
	switch command {
	case migrateUpCommand:
		migration.Up()
	case migrateDownCommand:
		migration.Down()
	case migrateNewCommand:
		if len(args) < 3 {
			log.Fatalf("[MIGRATION] Require a file name on %s command", migrateNewCommand)
		}
		fileName := args[2]
		migration.New(fileName)
	default:
		log.Fatalf("[MIGRATION] Invalid command")
	}
}
