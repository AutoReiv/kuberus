package handlers

import (
	"context"
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"

	"github.com/coreos/go-oidc"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

var (
	oidcProvider *oidc.Provider
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
)

// SetOIDCConfigHandler allows an admin to set the OIDC configuration.
func SetOIDCConfigHandler(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}

	username, _ := c.Get("username").(string)
	isAdmin, ok := c.Get("isAdmin").(bool)
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "Unable to determine admin status")
	}

	if !isAdmin {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Forbidden"})
	}

	var config auth.OIDCConfig
	if err := c.Bind(&config); err != nil {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Invalid request payload", err, "Failed to bind OIDC config request")
	}

	if err := auth.SetOIDCConfig(&config); err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to set OIDC configuration", err, "Failed to set OIDC configuration")
	}

	initOIDCProvider(config)

	utils.Logger.Info("OIDC configuration set successfully", zap.String("username", username))
	utils.LogAuditEvent(c.Request(), "set_oidc_config", "OIDC", "N/A")
	return c.JSON(http.StatusOK, map[string]string{"message": "OIDC configuration set successfully"})
}

// OIDCAuthHandler handles the OIDC authentication flow.
func OIDCAuthHandler(c echo.Context) error {
	if c.Request().Method != http.MethodGet {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}

	// Check if the OIDC configuration is set
	_, err := auth.GetOIDCConfig()
	if err != nil {
		utils.Logger.Error("OIDC configuration not set", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "OIDC configuration not set"})
	}

	state := "random" // You should generate a random state string for security
	authURL := oauth2Config.AuthCodeURL(state)
	utils.Logger.Info("Redirecting to OIDC provider for authentication", zap.String("authURL", authURL))
	return c.Redirect(http.StatusFound, authURL)
}

// OIDCCallbackHandler handles the OIDC callback.
func OIDCCallbackHandler(c echo.Context) error {
	if c.Request().Method != http.MethodGet {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}

	code := c.QueryParam("code")
	if code == "" {
		return utils.LogAndRespondError(c, http.StatusBadRequest, "Invalid request", nil, "Missing code parameter in OIDC callback")
	}

	utils.Logger.Info("Received authorization code", zap.String("code", code))

	oauth2Token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		utils.LogAuditEvent(c.Request(), "oidc_callback_failed", "N/A", "N/A")
		return utils.LogAndRespondError(c, http.StatusUnauthorized, "Authentication failed", err, "Failed to exchange token in OIDC callback")
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		utils.LogAuditEvent(c.Request(), "oidc_callback_failed", "N/A", "N/A")
		return utils.LogAndRespondError(c, http.StatusUnauthorized, "Authentication failed", nil, "Failed to retrieve ID token from OAuth2 token")
	}

	idToken, err := verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		utils.LogAuditEvent(c.Request(), "oidc_callback_failed", "N/A", "N/A")
		return utils.LogAndRespondError(c, http.StatusUnauthorized, "Authentication failed", err, "Failed to verify ID token in OIDC callback")
	}

	var claims struct {
		Email string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		utils.LogAuditEvent(c.Request(), "oidc_callback_failed", "N/A", "N/A")
		return utils.LogAndRespondError(c, http.StatusUnauthorized, "Authentication failed", err, "Failed to parse ID token claims in OIDC callback")
	}

	// Store OIDC user in the database if not already present
	if err := auth.CreateUserIfNotExists(claims.Email, "oidc"); err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to store user", err, "Failed to store OIDC user in database")
	}

	// Check if the user is an admin
	isAdmin := auth.IsAdmin(claims.Email)

	// Generate JWT token
	token, err := auth.GenerateJWT(claims.Email, isAdmin)
	if err != nil {
		return utils.LogAndRespondError(c, http.StatusInternalServerError, "Failed to generate token", err, "Failed to generate JWT token in OIDC callback")
	}

	utils.Logger.Info("OIDC callback successful", zap.String("email", claims.Email))
	utils.LogAuditEvent(c.Request(), "oidc_callback_success", claims.Email, "N/A")
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

// initOIDCProvider initializes the OIDC provider and verifier.
func initOIDCProvider(config auth.OIDCConfig) {
	var err error
	oidcProvider, err = oidc.NewProvider(context.Background(), config.IssuerURL)
	if err != nil {
		utils.Logger.Fatal("Error creating OIDC provider", zap.Error(err))
	}

	oauth2Config = &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint:     oidcProvider.Endpoint(),
		RedirectURL:  config.CallbackURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier = oidcProvider.Verifier(&oidc.Config{ClientID: config.ClientID})
	utils.Logger.Info("OIDC provider initialized", zap.String("issuerURL", config.IssuerURL))
}
