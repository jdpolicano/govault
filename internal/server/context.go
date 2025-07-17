package server

import (
	"log"
	"os"
	"time"

	"github.com/jdpolicano/govault/internal/store"
)

type ContextConfig struct {
	DefaultTTL time.Duration
	SaltSize   int
	VaultPath  string
}

func DefaultConfig() ContextConfig {
	return ContextConfig{
		DefaultTTL: time.Hour * 24,
		SaltSize:   16,
		VaultPath:  "./.govault",
	}
}

type Context struct {
	Sessions *SessionMap
	Store    store.Store
	Config   ContextConfig
	Log      *log.Logger
}

func NewContext(config ContextConfig) *Context {
	sessMap := NewSessionMap()
	store := store.NewJSONStore(config.VaultPath)
	logger := log.New(os.Stdout, "server: ", log.Ldate|log.Ltime)
	return &Context{sessMap, store, config, logger}
}
