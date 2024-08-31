package handlers

import (
	"net/http"
	"rbac/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/openidConnect"
)

// OIDCConfig represents the OIDC configuration.
type OIDCConfig struct {
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
	IssuerURL    string `json:"issuer_url" binding:"required"`
	CallbackURL  string `json:"callback_url" binding:"required"`
}

// SetOIDCConfigHandler allows an admin to set the OIDC configuration.
func SetOIDCConfigHandler(c *gin.Context) {
	// Ensure the user is an admin
	username, exists := c.Get("username")
	if !exists || !auth.IsAdmin(username.(string)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	var config OIDCConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Store the OIDC configuration (in-memory for simplicity)
	auth.Mu.Lock()
	auth.Config = &auth.OIDCConfig{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		IssuerURL:    config.IssuerURL,
		CallbackURL:  config.CallbackURL,
	}
	auth.Mu.Unlock()

	provider, err := openidConnect.New(config.ClientID, config.ClientSecret, config.CallbackURL, config.IssuerURL, "openid-connect")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OIDC provider: " + err.Error()})
		return
	}
	goth.UseProviders(provider)
}

// OIDCAuthHandler handles the OIDC authentication flow.
func OIDCAuthHandler(c *gin.Context) {
	// Check if the OIDC configuration is set
	if auth.Config == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OIDC configuration not set"})
		return
	}

	// Set the provider name to "openid-connect"
	query := c.Request.URL.Query()
	query.Set("provider", "openid-connect")
	c.Request.URL.RawQuery = query.Encode()

	// Begin the authentication handler
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// OIDCCallbackHandler handles the OIDC callback.
func OIDCCallbackHandler(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to complete OIDC authentication: " + err.Error()})
		return
	}

	// Here you can create a session for the user or generate a JWT token
	token, err := auth.GenerateJWT(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
