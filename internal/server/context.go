package server

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	e "github.com/jdpolicano/govault/internal/server/errors"
	"github.com/jdpolicano/govault/internal/store"
	"github.com/jdpolicano/govault/internal/vault"
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

// CreateUserSession encapsulates session setup logic.
func (ctx *Context) CreateUserSession(username string, key *vault.Key) (string, error) {
	sessId, err := GenerateSessionID()
	if err != nil {
		return "", err
	}
	ctx.Log.Printf("starting session id=\"%s\"", sessId)
	sess := NewSession(username, key.AES, ctx.Config.DefaultTTL)
	ctx.Sessions.Set(sessId, sess)
	return sessId, nil
}

// CreateUserSession encapsulates session setup logic.
func (ctx *Context) ValidateTokenHeader(req *http.Request) (Session, error) {
	// check that the header exists
	auth := req.Header.Get("Authorization")
	var none Session
	if len(auth) == 0 {
		return none, e.MissingAuthorizationHeader
	}

	// check that the prefix is correct
	toke, found := strings.CutPrefix(auth, "Bearer govault-")
	if !found {
		return none, e.MalformedAuthorizationHeader
	}
	// check if the token has expired...
	sess, found := ctx.Sessions.Get(toke)
	if sess.Expired() {
		ctx.Sessions.Delete(toke)
		return sess, e.AuthorizationExpired
	}

	return sess, nil
}
