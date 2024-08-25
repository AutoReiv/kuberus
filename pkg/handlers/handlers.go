package handlers

import (
	"sync"
	"time"
)

// Shared variables and synchronization primitives
var (
	users         = make(map[string]string)
	adminExists   bool
	sessions      = make(map[string]*Session)
	requestCounts = make(map[string]int)
	mu            sync.Mutex
)

// Session represents a user session
type Session struct {
	Username string
	Token    string
	ExpireAt time.Time
}

// rateLimit implements rate limiting logic
func rateLimit(ip string) bool {
	mu.Lock()
	defer mu.Unlock()

	if requestCounts[ip] >= 5 {
		return false
	}
	requestCounts[ip]++
	return true
}
