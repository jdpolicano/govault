package login

import (
	"bytes"
	"net/http"

	"github.com/jdpolicano/govault/internal/server"
	"github.com/jdpolicano/govault/internal/vault"
)

// HTTP handler function for logging in and getting a new token.
// todo: we should be validating the request type is a post request.
func Handler(ctx *server.Context) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// verify request body is well formed
		r, err := server.ValidateAuthRequest(req)
		if err != nil {
			server.JSONResponse(w, server.NewInvalidBodyError())
			return
		}
		username, password := r.Username, r.Password

		// now that we have the user and password, lets check if we have a user by this name.
		record, exists := ctx.Store.GetUserInfo(username)
		if !exists {
			server.JSONResponse(w, server.NewNoSuchUserError(username))
			return
		}

		// recompute the keys from the user's password and the stored salt
		key, err := vault.NewKeyWithSalt(password, record.Salt)
		if err != nil {
			server.JSONResponse(w, server.NewServerError(err))
			return
		}

		// if they are not the same the password is wrong...
		if !bytes.Equal(key.Login, record.Login) {
			server.JSONResponse(w, server.NewCredentialError())
			return
		}

		// if they are the same, create a new session with the aes key in memory and return
		// a token to the user for future requests.
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
