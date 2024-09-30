package server

import (
	"rbac/pkg/handlers"
	"rbac/pkg/middleware"

	"github.com/labstack/echo/v4"
)

func registerAuditRoutes(e *echo.Echo, config *Config) {
	api := e.Group("/api")
	api.Use(middleware.AuthAndRBACMiddleware("", config.IsDevMode))

	api.GET("/audit-logs", middleware.AuthAndRBACMiddleware("view_audit_logs", config.IsDevMode)(handlers.GetAuditLogsHandler))
}
