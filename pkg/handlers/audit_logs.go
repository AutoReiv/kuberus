package handlers

import (
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/db"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
)

// AuditLog represents a single audit log entry.
type AuditLog struct {
	ID           int    `json:"id"`
	Action       string `json:"action"`
	ResourceName string `json:"resource_name"`
	Namespace    string `json:"namespace"`
	Timestamp    string `json:"timestamp"`
	Hash         string `json:"hash"`
}

// GetAuditLogsHandler handles the retrieval of audit logs.
func GetAuditLogsHandler(c echo.Context) error {
	username := c.Get("username").(string)
	isAdmin, ok := c.Get("isAdmin").(bool)
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "Unable to determine admin status")
	}

	if (!isAdmin && !auth.HasPermission(username, "view_audit_logs")) {
		return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to view audit logs")
	}

	rows, err := db.DB.Query("SELECT action, resource_name, namespace, timestamp, hash FROM audit_logs ORDER BY timestamp DESC")
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error retrieving audit logs", err, "Failed to retrieve audit logs")
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		if err := rows.Scan(&log.Action, &log.ResourceName, &log.Namespace, &log.Timestamp, &log.Hash); err != nil {
			return utils.LogAndRespondError(c, http.StatusInternalServerError, "Error scanning audit log", err, "Failed to scan audit log")
		}
		logs = append(logs, log)
	}

	return c.JSON(http.StatusOK, logs)
}