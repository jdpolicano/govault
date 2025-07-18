package set

import (
	"encoding/json"
	"net/http"

	"github.com/jdpolicano/govault/internal/server"
	e "github.com/jdpolicano/govault/internal/server/errors"
	"github.com/jdpolicano/govault/internal/store"
	"github.com/jdpolicano/govault/internal/vault"
)

type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// HTTP handler function for logging in and getting a new token.
// todo: we should be validating the request type is a post request.
func Handler(ctx *server.Context) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// verify request has a correct header
		sess, err := ctx.ValidateTokenHeader(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// okay we have a valid session, lets validate the body of the request
		var body SetRequest
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&body); err != nil {
			http.Error(w, e.InvalidRequestBody.Error(), http.StatusBadRequest)
			return
		}

		// Now that we have the body, lets set this key
		cipher, nonce, err := vault.Encrypt(sess.Key, body.Value)
		if err != nil {
			ctx.Log.Printf("err encrypting key %s", err)
			http.Error(w, e.UnexpectedServerError.Error(), http.StatusInternalServerError)
			return
		}

		if err := setKey(ctx.Store, sess.User, body.Key, cipher, nonce); err != nil {
			ctx.Log.Printf("err setting key %s", err)
			http.Error(w, e.UnexpectedServerError.Error(), http.StatusInternalServerError)
			return
		}

		server.JSONResponse(w, server.NewResponse(http.StatusOK, "OK", nil))
	}
}

func setKey(s store.Store, user, key string, cipher, nonce []byte) error {
	return s.Set(user, key, store.CipherText{Nonce: nonce, Text: cipher})
}
