package routes

import (
	"net/http"

	"github.com/jdpolicano/govault/internal/server"
	"github.com/jdpolicano/govault/internal/store"
	"github.com/jdpolicano/govault/internal/vault"
)

// GetCreateRoute returns an HTTP handler function for creating a new user and session.
func GetRegisterRoute(ctx *server.Context) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		reqData, err := validateAuthRequest(req)
		if err != nil {
			server.JSONResponse(w, server.NewInvalidBodyError())
			return
		}

		username, password := reqData.Username, reqData.Password

		if ctx.Store.HasUser(username) {
			ctx.Log.Printf("error creating user \"%s\" already exists", username)
			server.JSONResponse(w, server.NewServerError(store.UserAlreadyExistsError{}))
			return
		}

		key, err := vault.NewKey(password, ctx.Config.SaltSize)
		if err != nil {
			ctx.Log.Printf("error creating user keys %s %v", username, err)
			server.JSONResponse(w, server.NewServerError(err))
			return
		}

		ctx.Log.Printf("successfully derived salt for user \"%s\"", username)
		login, salt := key.Base64LoginKey(), key.Base64Salt()

		if err = ctx.Store.AddUser(username, login, salt); err != nil {
			ctx.Log.Printf("error adding user \"%s\" %v", username, err)
			routeStoreError(w, err)
			return
		}

		ctx.Log.Printf("successfully added user \"%s\"", username)

		token, err := createUserSession(ctx, username, key)
		if err != nil {
			ctx.Log.Printf("error creating session for user \"%s\" %s", username, err)
			server.JSONResponse(w, server.NewServerError(err))
			return
		}

		sendToken(w, token)
		ctx.Log.Println("response sent")
	}
}

// sendToken sends a successful response with the session token.
func sendToken(w http.ResponseWriter, token string) {
	res := TokenSuccess{Token: token}
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
