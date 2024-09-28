package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	createTables()
	seedInitialData() // Call the seed function here
}

func createTables() {
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password TEXT NOT NULL,
		source TEXT NOT NULL,
		is_admin BOOLEAN NOT NULL DEFAULT 0
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
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
        hash TEXT
    );`

	createCertTable := `
	CREATE TABLE IF NOT EXISTS certificates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		cert BLOB NOT NULL,
		key BLOB NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createRolesTable := `
	CREATE TABLE IF NOT EXISTS roles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);`

	createPermissionsTable := `
	CREATE TABLE IF NOT EXISTS permissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);`

	createRolePermissionsTable := `
	CREATE TABLE IF NOT EXISTS role_permissions (
		role_id INTEGER,
		permission_id INTEGER,
		FOREIGN KEY (role_id) REFERENCES roles(id),
		FOREIGN KEY (permission_id) REFERENCES permissions(id)
	);`

	createUserRolesTable := `
	CREATE TABLE IF NOT EXISTS user_roles (
		user_id TEXT,
		role_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(username),
		FOREIGN KEY (role_id) REFERENCES roles(id)
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

	_, err = DB.Exec(createCertTable)
	if err != nil {
		log.Fatalf("Error creating certificates table: %v", err)
	}

	_, err = DB.Exec(createRolesTable)
	if err != nil {
		log.Fatalf("Error creating roles table: %v", err)
	}

	_, err = DB.Exec(createPermissionsTable)
	if err != nil {
		log.Fatalf("Error creating permissions table: %v", err)
	}

	_, err = DB.Exec(createRolePermissionsTable)
	if err != nil {
		log.Fatalf("Error creating role_permissions table: %v", err)
	}

	_, err = DB.Exec(createUserRolesTable)
	if err != nil {
		log.Fatalf("Error creating user_roles table: %v", err)
	}
}
func seedInitialData() {
	roles := []string{"admin", "editor", "viewer"}
	permissions := []string{
		"manage_namespaces", "list_namespaces", "create_namespace", "delete_namespace",
		"list_roles", "create_role", "update_role", "delete_role", "view_role_details",
		"list_rolebindings", "create_rolebinding", "update_rolebinding", "delete_rolebinding", "view_rolebinding_details",
		"list_clusterroles", "create_clusterrole", "update_clusterrole", "delete_clusterrole", "view_clusterrole_details",
		"list_clusterrolebindings", "create_clusterrolebinding", "update_clusterrolebinding", "delete_clusterrolebinding", "view_clusterrolebinding_details",
		"list_resources",
		"list_serviceaccounts", "create_serviceaccount", "delete_serviceaccount", "view_serviceaccount_details",
		"list_users", "view_user_details",
		"list_groups", "view_group_details",
		"view_user_roles",
		"view_audit_logs",
		"simulate",
	}

	for _, role := range roles {
		_, err := DB.Exec("INSERT OR IGNORE INTO roles (name) VALUES (?)", role)
		if err != nil {
			log.Fatalf("Error seeding roles: %v", err)
		}
	}

	for _, permission := range permissions {
		_, err := DB.Exec("INSERT OR IGNORE INTO permissions (name) VALUES (?)", permission)
		if err != nil {
			log.Fatalf("Error seeding permissions: %v", err)
		}
	}

	// Assign permissions to roles (example: admin gets all permissions)
	for _, permission := range permissions {
		var permissionID int
		err := DB.QueryRow("SELECT id FROM permissions WHERE name = ?", permission).Scan(&permissionID)
		if err != nil {
			log.Fatalf("Error fetching permission ID: %v", err)
		}

		_, err = DB.Exec("INSERT OR IGNORE INTO role_permissions (role_id, permission_id) VALUES ((SELECT id FROM roles WHERE name = 'admin'), ?)", permissionID)
		if err != nil {
			log.Fatalf("Error assigning permissions to admin role: %v", err)
		}
	}
}