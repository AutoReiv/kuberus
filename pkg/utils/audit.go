package utils

import (
    "log"
    "rbac/pkg/db"
)

// LogAuditEvent logs an audit event for role and role binding changes.
func LogAuditEvent(action, resourceName, namespace string) {
	log.Printf("Audit event: action=%s, resource=%s, namespace=%s", action, resourceName, namespace)

	_, err := db.DB.Exec("INSERT INTO audit_logs (action, resource_name, namespace) VALUES (?, ?, ?)", action, resourceName, namespace)
    if err != nil {
        log.Printf("Error logging audit event to database: %v", err)
    }
}
