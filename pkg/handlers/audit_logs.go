package handlers

import (
	"net/http"
	"rbac/pkg/db"

	"github.com/labstack/echo/v4"
)

// AuditLog represents a single audit log entry.
type AuditLog struct {
	ID           int    `json:"id"`
	Action       string `json:"action"`
	ResourceName string `json:"resource_name"`
	Namespace    string `json:"namespace"`
	Timestamp    string `json:"timestamp"`
}

// GetAuditLogsHandler handles the retrieval of audit logs.
func GetAuditLogsHandler(c echo.Context) error {
	rows, err := db.DB.Query("SELECT username, action, resource_name, namespace, ip_address, timestamp, hash FROM audit_logs ORDER BY timestamp DESC")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve audit logs: " + err.Error()})
	}
	defer rows.Close()

	var logs []map[string]string
	for rows.Next() {
		var username, action, resourceName, namespace, ipAddress, timestamp, hash string
		if err := rows.Scan(&username, &action, &resourceName, &namespace, &ipAddress, &timestamp, &hash); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to scan audit log: " + err.Error()})
		}
		logs = append(logs, map[string]string{
			"username":     username,
			"action":       action,
			"resourceName": resourceName,
			"namespace":    namespace,
			"ipAddress":    ipAddress,
			"timestamp":    timestamp,
			"hash":         hash,
		})
	}

	return c.JSON(http.StatusOK, logs)
}
