package handlers

import (
    "net/http"
    "github.com/labstack/echo/v4"
)

// TokenHandler handles requests for retrieving the token from localStorage.
func TokenHandler(c echo.Context) error {
    token := c.QueryParam("token")
    if token == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Token is required"})
    }
    return c.JSON(http.StatusOK, map[string]string{"token": token})
}
