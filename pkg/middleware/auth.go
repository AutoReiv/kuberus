package middleware

import (
	"net/http"
	"rbac/pkg/auth"
	"strings"

	"github.com/labstack/echo/v4"
)

// AuthAndRBACMiddleware validates the JWT token and checks if the user has the required permissions.
func AuthAndRBACMiddleware(permission string, isDevMode bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if isDevMode {
				return next(c)
			}

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := auth.ValidateJWT(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: "+err.Error())
			}

			// Add username and admin status to context
			c.Set("username", claims.Username)
			c.Set("isAdmin", claims.IsAdmin)

			// Check permissions
			if !auth.CheckPermission(claims.Username, permission) {
				return auth.LogAndRespondPermissionDenied(c, claims.Username, permission)
			}

			return next(c)
		}
	}
}
