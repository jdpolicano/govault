package register

import (
	"net/http"

	"github.com/jdpolicano/govault/internal/server"
	"github.com/jdpolicano/govault/internal/store"
	"github.com/jdpolicano/govault/internal/vault"
)

// HTTP handler function for creating a new user and session.
// todo: we should be validating the request type is a post request.
func Handler(ctx *server.Context) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// verify the reuqest body is well formed
		reqData, err := server.ValidateAuthRequest(req)
		if err != nil {
			server.JSONResponse(w, server.NewInvalidBodyError())
			return
		}
		username, password := reqData.Username, reqData.Password

		// check if this user already exists in the store
		if ctx.Store.HasUser(username) {
			ctx.Log.Printf("error creating user \"%s\" already exists", username)
			server.JSONResponse(w, server.NewServerError(store.UserAlreadyExistsError{}))
			return
		}

		// if not, then generate keys for this password and a new random salt for it.
		key, err := vault.NewKey(password, ctx.Config.SaltSize)
		if err != nil {
			ctx.Log.Printf("error creating user keys %s %v", username, err)
			server.JSONResponse(w, server.NewServerError(err))
			return
		}
		ctx.Log.Printf("successfully derived salt for user \"%s\"", username)

		// add the user to the store with the login key (for later authentication, NOT for encrypting/decrypting secrets)
		// and the salt that was used to derive that key
		if err = ctx.Store.AddUser(username, key.Login, key.Salt); err != nil {
			ctx.Log.Printf("error adding user \"%s\" %v", username, err)
			routeStoreError(w, err)
			return
		}
		ctx.Log.Printf("successfully added user \"%s\"", username)

		// issue a token to the user at this point so they won't need to call the login route separately.
		token, err := ctx.CreateUserSession(username, key)
		if err != nil {
			ctx.Log.Printf("error creating session for user \"%s\" %s", username, err)
			server.JSONResponse(w, server.NewServerError(err))
			return
		}

		server.SendToken(w, token)
		ctx.Log.Println("response sent")
	}
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
