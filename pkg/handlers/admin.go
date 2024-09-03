package handlers

import (
	"encoding/json"
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"
)

// CreateAdminRequest represents the request payload for creating an admin account.
type CreateAdminRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required,eqfield=Password"`
}

// CreateAdminHandler handles the creation of an admin account.
func CreateAdminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Sanitize user input
	username := utils.SanitizeInput(req.Username)
	password := utils.SanitizeInput(req.Password)

	// Validate password strength
	if !utils.IsStrongPassword(password) {
		http.Error(w, "Password does not meet strength requirements", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, "Error hashing password: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Acquire lock to synchronize access to shared data
	auth.Mu.Lock()
	defer auth.Mu.Unlock()

	// Check if an admin account already exists
	if auth.AdminExists {
		http.Error(w, "Admin account already exists", http.StatusConflict)
		return
	}

	// Store the admin account information
	auth.Users[username] = hashedPassword
	auth.AdminExists = true

	utils.WriteJSON(w, map[string]string{"message": "Admin account created successfully"})
}
