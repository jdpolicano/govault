package routes

import (
	"encoding/base64"
	"net/http"

	"github.com/jdpolicano/govault/internal/server"
	"github.com/jdpolicano/govault/internal/vault"
)

func GetLoginRoute(ctx *server.Context) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		r, err := validateAuthRequest(req)
		if err != nil {
			server.JSONResponse(w, server.NewInvalidBodyError())
			return
		}

		// now that we have the user and password, lets check if we have a user by this name.
		record, exists := ctx.Store.GetUserInfo(r.Username)
		if !exists {
			server.JSONResponse(w, server.NewNoSuchUserError(r.Username))
			return
		}

		key, err := vault.KeyFromSaltString(r.Password, record.Salt)
		if err != nil {
			server.JSONResponse(w, server.NewServerError(err))
			return
		}

		if ok := verifyKeyAgainstCache(record.Login, key); !ok {
			server.JSONResponse(w, server.NewCredentialError())
			return
		}

		token, err := createUserSession(ctx, r.Username, key)
		if err != nil {
			ctx.Log.Printf("error creating session for user \"%s\" %s", r.Username, err)
			server.JSONResponse(w, server.NewServerError(err))
			return
		}

		sendToken(w, token)
		ctx.Log.Println("response sent")
	}
}

func verifyKeyAgainstCache(cachedLogin string, key *vault.Key) bool {
	return cachedLogin == base64.RawStdEncoding.EncodeToString(key.Login)
}
