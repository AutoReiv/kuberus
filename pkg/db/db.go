package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	createTables()
}

func createTables() {
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password TEXT NOT NULL
	);`

	createOIDCTable := `
	CREATE TABLE IF NOT EXISTS oidc_config (
		id INTEGER PRIMARY KEY,
		client_id TEXT NOT NULL,
		client_secret TEXT NOT NULL,
		issuer_url TEXT NOT NULL,
		callback_url TEXT NOT NULL
	);`

	createAuditLogTable := `
    CREATE TABLE IF NOT EXISTS audit_logs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        action TEXT NOT NULL,
        resource_name TEXT NOT NULL,
        namespace TEXT NOT NULL,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err := DB.Exec(createUserTable)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	_, err = DB.Exec(createOIDCTable)
	if err != nil {
		log.Fatalf("Error creating oidc_config table: %v", err)
	}

	_, err = DB.Exec(createAuditLogTable)
	if err != nil {
		log.Fatalf("Error creating audit_logs table: %v", err)
	}
}
