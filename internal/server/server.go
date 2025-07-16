package server

import (
	"time"
)

type ServerConfig struct {
	defaultTTL time.Duration
	vaultPath  string
}

type Session struct {
	pw     string // the password of this session's user
	expiry int64  // unix timestamp of how long this session is valid for.
}

type Server struct {
	sessions   map[string]Session // an in memory map of the available sessions to the password associated with that session
	defaultTTL time.Duration
}

func NewServer(config *ServerConfig) *Server {
	return &Server{map[string]Session{}, config.defaultTTL}
}

func (s *Server) NewSession(pw string) (string, error) {
	ttl := time.Now().Add(s.defaultTTL).Unix()
	sess := Session{pw, ttl}
	key, err := GenerateSessionID()
	if err != nil {
		return "", err
	}
	s.sessions[key] = sess
	return key, err
}
