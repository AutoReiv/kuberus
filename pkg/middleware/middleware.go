package middleware

import (
	"net/http"
	"rbac/pkg/auth"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	jwtMiddleware "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ApplyMiddlewares applies all the middlewares to the given Echo instance.
func ApplyMiddlewares(e *echo.Echo, isDevMode bool) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
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

		// Add username to context
		c.Set("username", claims.Username)
		return next(c)
	}
}

// JWTMiddleware returns the JWT middleware configuration.
func JWTMiddleware() echo.MiddlewareFunc {
	return jwtMiddleware.WithConfig(jwtMiddleware.Config{
		SigningKey: auth.JwtKey,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.Claims)
		},
	})
}
