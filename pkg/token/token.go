package token

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

var (
	currentToken string
	mu           sync.Mutex
	// Token rotation interval
	expirationDate   time.Time
	rotationInterval = 30 * time.Second
)

const REQUESTED_BY = "WASM_Authors"

func init() {
	rotateToken() // Initial token generation
	// Set up a ticker to rotate the token periodically
	ticker := time.NewTicker(rotationInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				rotateToken()
			}
		}
	}()
}

// rotateToken generates a new token and stores it
func rotateToken() {
	mu.Lock()
	defer mu.Unlock()
	expirationDate = time.Now().Add(rotationInterval)
	currentToken = generateSecureToken()
}

// GenerateSecureToken creates a new secure token.
func generateSecureToken() string {
	b := make([]byte, 64) // 64 bytes = 512 bits
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // Handle error appropriately in production code
	}
	return base64.StdEncoding.EncodeToString(b)
}

// GetCurrentToken returns the current active token.
// This can be called from your WASM module to retrieve the current token.
func GetCurrentToken() (string, time.Time) {
	mu.Lock()
	defer mu.Unlock()
	return currentToken, expirationDate
}

// ValidateToken checks if the provided token matches the current active token.
func ValidateToken(token string) bool {
	mu.Lock()
	defer mu.Unlock()
	return token == currentToken
}
