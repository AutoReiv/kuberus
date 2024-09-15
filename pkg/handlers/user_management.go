package handlers

import (
	"encoding/json"
	"net/http"
	"rbac/pkg/auth"
	"rbac/pkg/utils"

	"k8s.io/client-go/kubernetes"
)

// CreateUserRequest represents the request payload for creating a user.
type CreateUserRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required,eqfield=Password"`
}

// UpdateUserRequest represents the request payload for updating a user's password.
type UpdateUserRequest struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required,eqfield=Password"`
}

// DeleteUserRequest represents the request payload for deleting a user.
type DeleteUserRequest struct {
	Username string `json:"username" binding:"required"`
}

// UserManagementHandler handles user management-related requests.
func UserManagementHandler(clientset *kubernetes.Clientset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateUser(w, r)
		case http.MethodPut:
			handleUpdateUser(w, r)
		case http.MethodDelete:
			handleDeleteUser(w, r)
		case http.MethodGet:
			handleListUsers(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
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

	// Create user
	if err := auth.CreateUser(username, password); err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogAuditEvent(r, "create_user", username, "N/A")
	utils.WriteJSON(w, map[string]string{"message": "User account created successfully"})
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var req UpdateUserRequest
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

	// Update user
	if err := auth.UpdateUser(username, password); err != nil {
		http.Error(w, "Error updating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogAuditEvent(r, "update_user", username, "N/A")
	utils.WriteJSON(w, map[string]string{"message": "User account updated successfully"})
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	var req DeleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Sanitize user input
	username := utils.SanitizeInput(req.Username)

	// Delete user
	if err := auth.DeleteUser(username); err != nil {
		http.Error(w, "Error deleting user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogAuditEvent(r, "delete_user", username, "N/A")
	utils.WriteJSON(w, map[string]string{"message": "User account deleted successfully"})
}

func handleListUsers(w http.ResponseWriter) {
	users, err := auth.GetAllUsers()
	if err != nil {
		http.Error(w, "Error retrieving users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, users)
}
