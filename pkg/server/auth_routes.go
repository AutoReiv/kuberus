package server

import (
	"rbac/pkg/handlers"

	"github.com/labstack/echo/v4"
)

func registerAuthRoutes(e *echo.Echo) {
	e.POST("/auth/login", handlers.LoginHandler)
	e.GET("/auth/oidc/login", handlers.OIDCAuthHandler)
	e.GET("/auth/oidc/callback", handlers.OIDCCallbackHandler)
}