package auth

import (
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var (
	Mu          sync.Mutex
	Users       = make(map[string]string)
	AdminExists bool
)

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
