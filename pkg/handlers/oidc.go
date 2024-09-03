package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"

	"github.com/coreos/go-oidc"
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
func SetOIDCConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ensure the user is an admin
	username := r.Context().Value("username").(string)
	if !auth.IsAdmin(username) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var config OIDCConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, "Failed to create OIDC provider: "+err.Error(), http.StatusInternalServerError)
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
	utils.WriteJSON(w, map[string]string{"message": "OIDC configuration set successfully"})
}

// OIDCAuthHandler handles the OIDC authentication flow.
func OIDCAuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the OIDC configuration is set
	if auth.Config == nil {
		http.Error(w, "OIDC configuration not set", http.StatusBadRequest)
		return
	}

	state := "random" // You should generate a random state string for security
	authURL := oauth2Config.AuthCodeURL(state)
	http.Redirect(w, r, authURL, http.StatusFound)
}

// OIDCCallbackHandler handles the OIDC callback.
func OIDCCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Log the authorization code for debugging
	log.Printf("Authorization code: %s", code)

	oauth2Token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Failed to exchange token: %v", err) // Log the error for debugging
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	idToken, err := verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	var claims struct {
		Email string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Here you can create a session for the user or generate a JWT token
	token, err := auth.GenerateJWT(claims.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"token": token})
}
