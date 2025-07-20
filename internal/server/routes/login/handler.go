package login

import (
	"bytes"
	"net/http"

	"github.com/jdpolicano/govault/internal/server"
	"github.com/jdpolicano/govault/internal/server/middleware"
	"github.com/jdpolicano/govault/internal/vault"
)

// HTTP handler function for logging in and getting a new token.
// todo: we should be validating the request type is a post request.
func Handler(refs *server.ServerRefs) http.HandlerFunc {

	handle := func(w http.ResponseWriter, req *http.Request) {
		body := req.Context().Value(server.BodyKey{}).(server.AuthCredentials)
		username, password := body.Username, body.Password
		// now that we have the user and password, lets check if we have a user by this name.
		record, exists := refs.Store.GetUserInfo(username)
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
		token, err := refs.Sessions.CreateUserSession(username, key.AES, refs.Config.DefaultTTL)
		if err != nil {
			refs.Log.Printf("error creating session for user \"%s\" %s", username, err)
			server.JSONResponse(w, server.NewServerError(err))
			return
		}
		server.SendToken(w, token)
		refs.Log.Println("response sent")
	}

	return middleware.Chain(handle,
		middleware.Logging(refs.Log),
		middleware.ParseJSONBody[server.AuthCredentials](),
	)

}
