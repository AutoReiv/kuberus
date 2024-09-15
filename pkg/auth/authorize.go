package auth

import (
	"rbac/pkg/db"
	"rbac/pkg/utils"

	"golang.org/x/crypto/bcrypt"
	"go.uber.org/zap"
)

// OIDCConfig represents the OIDC configuration.
type OIDCConfig struct {
	ClientID     string `json:"client_id" binding:"required"`
	ClientSecret string `json:"client_secret" binding:"required"`
	IssuerURL    string `json:"issuer_url" binding:"required"`
	CallbackURL  string `json:"callback_url" binding:"required"`
}

// HashPassword hashes a plain text password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with a plain text password.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// AuthenticateUser authenticates the user credentials.
func AuthenticateUser(username, password string) bool {
	var hashedPassword string
	err := db.DB.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		utils.Logger.Error("Error querying user password", zap.Error(err))
		return false
	}

	return CheckPasswordHash(password, hashedPassword)
}

// IsAdmin checks if a user is an admin.
func IsAdmin(username string) bool {
	var isAdmin bool
	err := db.DB.QueryRow("SELECT is_admin FROM users WHERE username = ?", username).Scan(&isAdmin)
	if err != nil {
		utils.Logger.Error("Error querying user admin status", zap.Error(err))
		return false
	}
	return isAdmin
}

// CreateUser creates a new user.
func CreateUser(username, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		utils.Logger.Error("Error hashing password", zap.Error(err))
		return err
	}

	_, err = db.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		utils.Logger.Error("Error creating user", zap.Error(err))
		return err
	}

	return nil
}

// UpdateUser updates an existing user's password.
func UpdateUser(username, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		utils.Logger.Error("Error hashing password", zap.Error(err))
		return err
	}

	_, err = db.DB.Exec("UPDATE users SET password = ? WHERE username = ?", hashedPassword, username)
	if err != nil {
		utils.Logger.Error("Error updating user password", zap.Error(err))
		return err
	}

	return nil
}

// DeleteUser deletes a user.
func DeleteUser(username string) error {
	_, err := db.DB.Exec("DELETE FROM users WHERE username = ?", username)
	if err != nil {
		utils.Logger.Error("Error deleting user", zap.Error(err))
		return err
	}

	return nil
}

// SetOIDCConfig sets the OIDC configuration.
func SetOIDCConfig(config *OIDCConfig) error {
	_, err := db.DB.Exec("INSERT INTO oidc_config (client_id, client_secret, issuer_url, callback_url) VALUES (?, ?, ?, ?)",
		config.ClientID, config.ClientSecret, config.IssuerURL, config.CallbackURL)
	if err != nil {
		utils.Logger.Error("Error setting OIDC configuration", zap.Error(err))
	}
	return err
}

// GetOIDCConfig retrieves the OIDC configuration.
func GetOIDCConfig() (*OIDCConfig, error) {
	var config OIDCConfig
	err := db.DB.QueryRow("SELECT client_id, client_secret, issuer_url, callback_url FROM oidc_config LIMIT 1").
		Scan(&config.ClientID, &config.ClientSecret, &config.IssuerURL, &config.CallbackURL)
	if err != nil {
		utils.Logger.Error("Error retrieving OIDC configuration", zap.Error(err))
		return nil, err
	}
	return &config, nil
}

// User represents a user account.
type User struct {
	Username string `json:"username"`
}

// GetAllUsers retrieves all user accounts.
func GetAllUsers() ([]User, error) {
	rows, err := db.DB.Query("SELECT username FROM users")
	if err != nil {
		utils.Logger.Error("Error retrieving users", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Username); err != nil {
			utils.Logger.Error("Error scanning user", zap.Error(err))
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}