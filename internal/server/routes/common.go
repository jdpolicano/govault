package routes

import (
	"encoding/json"
	"net/http"

	"github.com/jdpolicano/govault/internal/server"
	"github.com/jdpolicano/govault/internal/vault"
)

type TokenSuccess struct {
	Token string `json:"token"`
}

type AuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// validateRequest decodes and validates the request body.
func validateAuthRequest(req *http.Request) (AuthCredentials, error) {
	var reqData AuthCredentials
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqData)
	return reqData, err
}

// createUserSession encapsulates session setup logic.
func createUserSession(ctx *server.Context, username string, key *vault.Key) (string, error) {
	sessId, err := server.GenerateSessionID()
	if err != nil {
		return "", err
	}
	ctx.Log.Printf("starting session id=\"%s\"", sessId)
	sess := server.NewSession(username, key.AES, ctx.Config.DefaultTTL)
	ctx.Sessions.Set(sessId, sess)
	return sessId, nil
}
