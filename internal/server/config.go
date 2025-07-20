package server

import "time"

type ContextConfig struct {
	DefaultTTL time.Duration
	SaltSize   int
	VaultPath  string
}

func DefaultConfig() *ContextConfig {
	return &ContextConfig{
		DefaultTTL: time.Hour * 24,
		SaltSize:   16,
		VaultPath:  "./.govault",
	}
}
