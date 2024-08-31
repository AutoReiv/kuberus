package auth

import (
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var (
	Mu          sync.Mutex
	Users       = make(map[string]string)
	AdminExists bool
	Config      *OIDCConfig
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
	Mu.Lock()
	defer Mu.Unlock()

	hashedPassword, ok := Users[username]
	if !ok {
		return false
	}

	return CheckPasswordHash(password, hashedPassword)
}

// IsAdmin checks if a user is an admin.
func IsAdmin(username string) bool {
	Mu.Lock()
	defer Mu.Unlock()

	_, ok := Users[username]
	return ok && AdminExists
}
