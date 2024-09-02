package handlers

import (
	"context"
	"log"
	"net/http"
	"rbac/pkg/auth"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	oidcProvider *oidc.Provider
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
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

	// Log the OIDC configuration for debugging
	log.Printf("Setting OIDC Config: %+v", config)

	// Store the OIDC configuration (in-memory for simplicity)
	auth.Mu.Lock()
	auth.Config = &auth.OIDCConfig{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		IssuerURL:    config.IssuerURL,
		CallbackURL:  config.CallbackURL,
	}
	auth.Mu.Unlock()

	var err error
	oidcProvider, err = oidc.NewProvider(context.Background(), config.IssuerURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OIDC provider: " + err.Error()})
		return
	}

	oauth2Config = &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint:     oidcProvider.Endpoint(),
		RedirectURL:  config.CallbackURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier = oidcProvider.Verifier(&oidc.Config{ClientID: config.ClientID})
	c.JSON(http.StatusOK, gin.H{"message": "OIDC configuration set successfully"})
}

// OIDCAuthHandler handles the OIDC authentication flow.
func OIDCAuthHandler(c *gin.Context) {
	// Check if the OIDC configuration is set
	if auth.Config == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OIDC configuration not set"})
		return
	}

	state := "random" // You should generate a random state string for security
	authURL := oauth2Config.AuthCodeURL(state)
	c.Redirect(http.StatusFound, authURL)
}

// OIDCCallbackHandler handles the OIDC callback.
func OIDCCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	oauth2Token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to exchange token: " + err.Error()})
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No id_token field in oauth2 token"})
		return
	}

	idToken, err := verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify ID token: " + err.Error()})
		return
	}

	var claims struct {
		Email string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse ID token claims: " + err.Error()})
		return
	}

	// Here you can create a session for the user or generate a JWT token
	token, err := auth.GenerateJWT(claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
