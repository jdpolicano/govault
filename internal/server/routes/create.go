package routes

import (
	"encoding/json"
	"net/http"

	"github.com/jdpolicano/govault/internal/server"
	"github.com/jdpolicano/govault/internal/store"
	"github.com/jdpolicano/govault/internal/vault"
)

type CreateSuccess struct {
	Token string `json:"token"`
}

type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// GetCreateRoute returns an HTTP handler function for creating a new user and session.
func GetCreateRoute(ctx *server.Context) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		reqData, err := validateRequest(req)
		if err != nil {
			server.JSONResponse(w, server.NewInvalidBodyError())
			return
		}

		user, salt, err := vault.CreateUserFromPW(reqData.Username, reqData.Password)
		if err != nil {
			ctx.Log.Printf("error creating user %s: %v", reqData.Username, err)
			server.JSONResponse(w, server.NewServerError(err))
			return
		}
		ctx.Log.Printf("successfully derived salt for user \"%s\"", reqData.Username)

		if err = ctx.Store.AddUser(user); err != nil {
			ctx.Log.Printf("error adding user %s: %v", reqData.Username, err)
			routeStoreError(w, err)
			return
		}
		ctx.Log.Printf("successfully added user \"%s\"", reqData.Username)

		token, err := createUserSession(ctx, user, reqData.Password, salt)
		if err != nil {
			ctx.Log.Printf("error creating session for user %s: %v", reqData.Username, err)
			server.JSONResponse(w, server.NewServerError(err))
			return
		}

		sendToken(w, token)
		ctx.Log.Println("response sent")
	}
}

// validateRequest decodes and validates the request body.
func validateRequest(req *http.Request) (CreateRequest, error) {
	var reqData CreateRequest
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqData)
	return reqData, err
}

// createUserSession encapsulates session setup logic.
func createUserSession(ctx *server.Context, user store.User, password string, salt []byte) (string, error) {
	sessId, err := server.GenerateSessionID()
	if err != nil {
		return "", err
	}
	ctx.Log.Printf("starting session id=\"%s\"", sessId)

	aesKey, err := vault.CreateTextKeyFromPW("aes", password, salt)
	if err != nil {
		return "", err
	}

	sess := server.NewSession(user.Name, aesKey, ctx.Config.DefaultTTL)
	ctx.Sessions.Set(sessId, sess)
	return sessId, nil
}

// sendToken sends a successful response with the session token.
func sendToken(w http.ResponseWriter, token string) {
	res := CreateSuccess{Token: token}
	server.JSONResponse(w, server.NewServerSuccess(res))
}

// routeStoreError handles different store error types.
func routeStoreError(w http.ResponseWriter, err error) {
	switch err.(type) {
	case store.UserAlreadyExistsError:
		server.JSONResponse(w, server.NewClientError(err))
	default:
		server.JSONResponse(w, server.NewServerError(err))
	}
}
