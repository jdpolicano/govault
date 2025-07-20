package server

import (
	"log"
	"os"

	"github.com/jdpolicano/govault/internal/store"
)

type ServerRefs struct {
	Sessions *SessionMap
	Store    store.Store
	Config   *ContextConfig
	Log      *log.Logger
}

func NewServerRefs(config *ContextConfig) *ServerRefs {
	sessMap := NewSessionMap()
	store := store.NewJSONStore(config.VaultPath)
	logger := log.New(os.Stdout, "server: ", log.Ldate|log.Ltime)
	return &ServerRefs{sessMap, store, config, logger}
}
