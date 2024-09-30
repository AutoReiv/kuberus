package utils

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// LogAndRespondError logs the error and sends a JSON response with the error message.
func LogAndRespondError(c echo.Context, statusCode int, userMessage string, err error, logMessage string) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	Logger.Error(logMessage,
		zap.String("requestID", requestID),
		zap.String("user", c.Get("username").(string)),
		zap.Error(err),
	)
	return c.JSON(statusCode, map[string]string{"error": userMessage, "requestID": requestID})
}
