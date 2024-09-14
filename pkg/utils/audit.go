package utils

import (
    "log"
    "net/http"
    "rbac/pkg/db"
    "time"
)

// LogAuditEvent logs an audit event for role and role binding changes.
func LogAuditEvent(r *http.Request, action, resourceName, namespace string) {
    username := r.Context().Value("username").(string)
    ipAddress := r.RemoteAddr
    timestamp := time.Now().Format(time.RFC3339)

    log.Printf("Audit event: user=%s, action=%s, resource=%s, namespace=%s, ip=%s, timestamp=%s", username, action, resourceName, namespace, ipAddress, timestamp)

    _, err := db.DB.Exec("INSERT INTO audit_logs (username, action, resource_name, namespace, ip_address, timestamp) VALUES (?, ?, ?, ?, ?, ?)",
        username, action, resourceName, namespace, ipAddress, timestamp)
    if err != nil {
        log.Printf("Error logging audit event to database: %v", err)
    }
}