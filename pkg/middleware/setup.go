package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
}
