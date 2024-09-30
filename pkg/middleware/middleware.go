package middleware

import (
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// ApplyMiddlewares applies all the middlewares to the given Echo instance.
func ApplyMiddlewares(e *echo.Echo, isDevMode bool) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Secure headers middleware
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	}))

	// Apply rate limiting middleware
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))
}

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
			if !claims.IsAdmin && !auth.HasPermission(claims.Username, permission) {
				utils.Logger.Warn("Permission denied", zap.String("username", claims.Username), zap.String("permission", permission))
				return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to access this resource")
			}

			return next(c)
		}
	}
}
