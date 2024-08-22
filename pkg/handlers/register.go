package handlers

import (
	"encoding/json"
	"net/http"
	"rbac/pkg/utils"
	"sync"
)

// Mock user data store
var (
	users       = map[string]string{}
	adminExists = false
	mu          sync.Mutex
)

// RegisterRequest represents the registration request payload.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterResponse represents the registration response payload.
type RegisterResponse struct {
	Message string `json:"message"`
}

// RegisterHandler handles user registration.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		// Check if the admin account already exists
		if adminExists {
			http.Error(w, "Registration is closed", http.StatusForbidden)
			return
		}

		// Check if the username already exists
		if _, exists := users[req.Username]; exists {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}

		// Hash the password
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		// Store the user details
		users[req.Username] = hashedPassword

		// Set the adminExists flag to true
		adminExists = true

		resp := RegisterResponse{Message: "User registered successfully"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
