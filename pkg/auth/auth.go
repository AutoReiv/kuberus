package auth

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	Mu          sync.Mutex
	Users       = make(map[string]string)
	AdminExists bool
)

// Session represents a user session.
type Session struct {
	Username string
	Token    string
	ExpireAt time.Time
}

var sessions = make(map[string]*Session)

// StoreSession stores a new session.
func StoreSession(session *Session) {
	Mu.Lock()
	defer Mu.Unlock()
	sessions[session.Token] = session
}

// GetSession retrieves a session by token.
func GetSession(token string) (*Session, bool) {
	Mu.Lock()
	defer Mu.Unlock()
	session, exists := sessions[token]
	if !exists {
		return nil, false
	}
	// Check if the session has expired
	if session.ExpireAt.Before(time.Now()) {
		delete(sessions, token) // Clean up expired session
		return nil, false
	}
	return session, true
}

// GenerateSessionToken generates a secure random session token.
func GenerateSessionToken() string {
	return uuid.New().String()
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
