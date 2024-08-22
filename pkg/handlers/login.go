package handlers

import (
	"encoding/json"
	"net/http"
	"rbac/pkg/utils"
	"time"
)

// LoginRequest represents the login request payload.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response payload.
type LoginResponse struct {
	Message string `json:"message"`
}

// LoginHandler handles the login page.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Authenticate user
		hashedPassword, ok := users[req.Username]
		if !ok || !utils.CheckPasswordHash(req.Password, hashedPassword) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Set a cookie for session management
		expiration := time.Now().Add(24 * time.Hour)
		cookie := http.Cookie{Name: "session_token", Value: "authenticated", Expires: expiration}
		http.SetCookie(w, &cookie)

		resp := LoginResponse{Message: "Login successful"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
