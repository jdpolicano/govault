package routes

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/jdpolicano/govault/internal/server"
)

type LoginRespone struct {
	sessionKey string `json:sessionKey`
}
type LoginRequest struct {
	user string `json:user`
	password string `json:password`
}

func GetLoginRoute(s *server.Server) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		req.Body.
	}
}

// GenerateSessionID generates a cryptographically secure, random session ID.
func generateSessionID() (string, error) {
	// A common practice is to use a 32-byte (256-bit) random value for session IDs.
	// This provides sufficient entropy for security.
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}
	// Encode the random bytes to a URL-safe base64 string.
	// This makes the ID suitable for use in cookies or URLs.
	return base64.URLEncoding.EncodeToString(b), nil
}
