package server

import (
	"encoding/base64"
	"sync"
	"time"

	"github.com/jdpolicano/govault/internal/vault"
)

type Session struct {
	User string
	Key  string
	TTL  int64
}

func NewSession(user, key string, ttl time.Duration) Session {
	eol := time.Now().Add(ttl).Unix()
	return Session{user, key, eol}
}

type SessionMap struct {
	sync.RWMutex
	sessions map[string]Session // a map from a session key to
}

func NewSessionMap() *SessionMap {
	return &SessionMap{sessions: make(map[string]Session, 1024)}
}

func (s *SessionMap) Get(key string) (Session, bool) {
	s.RLock()
	defer s.RUnlock()
	sess, exists := s.sessions[key]
	return sess, exists
}

func (s *SessionMap) Set(key string, sess Session) {
	s.Lock()
	defer s.Unlock()
	s.sessions[key] = sess
}

// GenerateSessionID generates a cryptographically secure, random session ID.
func GenerateSessionID() (string, error) {
	// A common practice is to use a 32-byte (256-bit) random value for session IDs.
	// This provides sufficient entropy for security.
	b, err := vault.GenerateRandBytes(32)
	if err != nil {
		return "", err
	}
	// Encode the random bytes to a URL-safe base64 string.
	// This makes the ID suitable for use in cookies or URLs.
	return base64.URLEncoding.EncodeToString(b), nil
}
