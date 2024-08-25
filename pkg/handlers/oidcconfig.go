package handlers

import (
	"net/http"
	"os"

	"rbac/pkg/auth"

	"github.com/gin-gonic/gin"
)

func SetupOIDCConfigHandler(c *gin.Context) {
	// Parse the OIDC configuration from the request body
	var oidcConfig struct {
		ClientID     string `json:"clientID"`
		ClientSecret string `json:"clientSecret"`
		CallbackURL  string `json:"callbackURL"`
		Endpoint     string `json:"endpoint"`
	}
	if err := c.ShouldBindJSON(&oidcConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the OIDC environment variables
	os.Setenv("OIDC_CLIENT_ID", oidcConfig.ClientID)
	os.Setenv("OIDC_CLIENT_SECRET", oidcConfig.ClientSecret)
	os.Setenv("OIDC_CALLBACK_URL", oidcConfig.CallbackURL)
	os.Setenv("OIDC_ENDPOINT", oidcConfig.Endpoint)

	// Configure the OIDC provider
	auth.ConfigureOIDCProvider()

	c.JSON(http.StatusOK, gin.H{"message": "OIDC configuration updated successfully"})
}
