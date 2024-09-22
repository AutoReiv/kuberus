package utils

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// LogAndRespondError logs an error with context and sends a JSON response with the error message.
func LogAndRespondError(c echo.Context, statusCode int, userMessage string, err error, contextMessage string) error {
	Logger.Error(contextMessage, zap.Error(err))
	return c.JSON(statusCode, map[string]string{"error": userMessage})
}
