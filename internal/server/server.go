package server

import (
	"time"
)

type ServerConfig struct {
	defaultTTL time.Duration
	vaultPath  string
}

type Server struct {
	sessions *SessionMap
	config   ServerConfig
}

func NewServer(config ServerConfig) *Server {
	return &Server{NewSessionMap(), config}
}
