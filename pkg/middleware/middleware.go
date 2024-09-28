package middleware

import (
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	jwtMiddleware "github.com/labstack/echo-jwt/v4"
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

// AuthMiddleware validates the JWT token and sets the user information in the request context.
func AuthMiddleware(next echo.HandlerFunc, isDevMode bool) echo.HandlerFunc {
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
		return next(c)
	}
}

// RBACMiddleware checks if the user has the required permissions for the request.
func RBACMiddleware(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			username, ok := c.Get("username").(string)
			if !ok || username == "" {
				utils.Logger.Warn("Username not found in context")
				return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to access this resource with permission: "+permission)
			}

			isAdmin, ok := c.Get("isAdmin").(bool)
			if ok && isAdmin {
				// If the user is an admin, allow access to all resources
				utils.Logger.Info("Admin access granted", zap.String("username", username))
				return next(c)
			}

			utils.Logger.Info("Checking permission", zap.String("username", username), zap.String("permission", permission))
			if !auth.HasPermission(username, permission) {
				utils.Logger.Warn("Permission denied", zap.String("username", username), zap.String("permission", permission))
				return echo.NewHTTPError(http.StatusForbidden, "You do not have permission to access this resource with permission: "+permission+" for user: "+username)
			}
			return next(c)
		}
	}
}

// JWTMiddleware returns the JWT middleware configuration.
func JWTMiddleware() echo.MiddlewareFunc {
	return jwtMiddleware.WithConfig(jwtMiddleware.Config{
		SigningKey: auth.JwtKey,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.Claims)
		},
		SuccessHandler: func(c echo.Context) {
			user := c.Get("user")
			token, ok := user.(*jwt.Token)
			if !ok {
				utils.Logger.Error("Failed to assert user to *jwt.Token")
				return
			}

			claims, ok := token.Claims.(*auth.Claims)
			if !ok {
				utils.Logger.Error("Failed to assert token claims to *auth.Claims")
				return
			}

			c.Set("username", claims.Username)
			c.Set("isAdmin", claims.IsAdmin)
		},
	})

}