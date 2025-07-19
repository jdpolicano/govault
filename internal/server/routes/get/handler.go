package get

import (
	"encoding/json"
	"net/http"

	"github.com/jdpolicano/govault/internal/server"
	e "github.com/jdpolicano/govault/internal/server/errors"
	"github.com/jdpolicano/govault/internal/vault"
)

type GetRequest struct {
	Key string `json:"key"`
}

func Handler(ctx *server.Context) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		sess, err := ctx.ValidateTokenHeader(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var body GetRequest
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&body); err != nil {
			http.Error(w, e.InvalidRequestBody.Error(), http.StatusBadRequest)
			return
		}

		cipher, exists := ctx.Store.Get(sess.User, body.Key)
		if !exists {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		plain, err := vault.Decrypt(cipher.Nonce, sess.Key, cipher.Text)
		if err != nil {
			ctx.Log.Printf("err decrypting key %s", err)
			http.Error(w, e.UnexpectedServerError.Error(), http.StatusInternalServerError)
			return
		}

		server.JSONResponse(w, server.NewResponse(http.StatusOK, string(plain), nil))
	}
}
