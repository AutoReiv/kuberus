package handlers

import (
	"net/http"
	"rbac/pkg/db"
	"rbac/pkg/utils"
)

// AuditLog represents a single audit log entry.
type AuditLog struct {
	ID           int    `json:"id"`
	Action       string `json:"action"`
	ResourceName string `json:"resource_name"`
	Namespace    string `json:"namespace"`
	Timestamp    string `json:"timestamp"`
}

// GetAuditLogsHandler handles requests to retrieve audit logs.
func GetAuditLogsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, action, resource_name, namespace, timestamp FROM audit_logs ORDER BY timestamp DESC")
	if err != nil {
		http.Error(w, "Error retrieving audit logs: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		if err := rows.Scan(&log.ID, &log.Action, &log.ResourceName, &log.Namespace, &log.Timestamp); err != nil {
			http.Error(w, "Error scanning audit log: "+err.Error(), http.StatusInternalServerError)
			return
		}
		logs = append(logs, log)
	}

	utils.WriteJSON(w, logs)
}
