package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"rbac/pkg/db"
	"time"
)

// LogAuditEvent logs an audit event for role and role binding changes.
func LogAuditEvent(r *http.Request, action, resourceName, namespace string) {
	username, ok := r.Context().Value("username").(string)
	if !ok {
		log.Printf("Error: username not found in context")
		return
	}
	ipAddress := r.RemoteAddr
	timestamp := time.Now().Format(time.RFC3339)

	// Create a hash of the log entry
	logEntry := username + action + resourceName + namespace + ipAddress + timestamp
	hash := sha256.Sum256([]byte(logEntry))
	hashString := hex.EncodeToString(hash[:])

	log.Printf("Audit event: user=%s, action=%s, resource=%s, namespace=%s, ip=%s, timestamp=%s, hash=%s", username, action, resourceName, namespace, ipAddress, timestamp, hashString)

	_, err := db.DB.Exec("INSERT INTO audit_logs (username, action, resource_name, namespace, ip_address, timestamp, hash) VALUES (?, ?, ?, ?, ?, ?, ?)",
		username, action, resourceName, namespace, ipAddress, timestamp, hashString)
	if err != nil {
		log.Printf("Error logging audit event to database: %v", err)
	}
}
