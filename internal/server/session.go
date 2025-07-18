package server

import (
	"encoding/base64"
	"sync"
	"time"

	"github.com/jdpolicano/govault/internal/vault"
)

type Session struct {
	User string
	Key  []byte
	TTL  int64
}

func NewSession(user string, key []byte, ttl time.Duration) Session {
	eol := time.Now().Add(ttl).Unix()
	return Session{user, key, eol}
}

func (s Session) Expired() bool {
	expiry := time.Unix(s.TTL, 0)
	return time.Now().After(expiry)
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

func (s *SessionMap) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.sessions, key)
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
	return base64.RawStdEncoding.EncodeToString(b), nil
}
