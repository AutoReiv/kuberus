package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func OAuthLoginHandler(c *gin.Context) {
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func OAuthCallbackHandler(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Handle the authenticated user (e.g., create a session, store user information)
	// ...

	c.JSON(http.StatusOK, user)
}
