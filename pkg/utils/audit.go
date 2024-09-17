package utils

import (
    "crypto/sha256"
    "encoding/hex"
    "net/http"
    "rbac/pkg/db"
    "time"

    "go.uber.org/zap"
)

// LogAuditEvent logs an audit event for role and role binding changes.
func LogAuditEvent(r *http.Request, action, resourceName, namespace string) {
    timestamp := time.Now().Format(time.RFC3339)

    // Create a hash of the log entry
    logEntry := action + resourceName + namespace + timestamp
    hash := sha256.Sum256([]byte(logEntry))
    hashString := hex.EncodeToString(hash[:])

    zap.L().Info("Audit event",
        zap.String("action", action),
        zap.String("resource", resourceName),
        zap.String("namespace", namespace),
        zap.String("timestamp", timestamp),
        zap.String("hash", hashString),
    )

    _, err := db.DB.Exec("INSERT INTO audit_logs (action, resource_name, namespace, timestamp, hash) VALUES (?, ?, ?, ?, ?)",
        action, resourceName, namespace, timestamp, hashString)
    if err != nil {
        zap.L().Error("Error logging audit event to database", zap.Error(err))
    }
}
