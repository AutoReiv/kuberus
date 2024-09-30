package utils

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var Logger *zap.Logger

// InitLogger initializes the logger.
func InitLogger() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
}

// LogAndRespondError logs the error with additional context and sends a JSON response.
func LogAndRespondError(c echo.Context, statusCode int, userMessage string, err error, logMessage string) error {
	requestID := c.Response().Header().Get(echo.HeaderXRequestID)
	Logger.Error(logMessage,
		zap.String("requestID", requestID),
		zap.String("user", c.Get("username").(string)),
		zap.Error(err),
	)
	return c.JSON(statusCode, map[string]string{"error": userMessage, "requestID": requestID})
}
