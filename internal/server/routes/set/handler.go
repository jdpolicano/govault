package set

import (
	"net/http"

	"github.com/jdpolicano/govault/internal/server"
	e "github.com/jdpolicano/govault/internal/server/errors"
	"github.com/jdpolicano/govault/internal/server/middleware"
	"github.com/jdpolicano/govault/internal/store"
	"github.com/jdpolicano/govault/internal/vault"
)

type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// HTTP handler function for logging in and getting a new token.
// todo: we should be validating the request type is a post request.
func Handler(ctx *server.Context) http.HandlerFunc {
	handle := func(w http.ResponseWriter, req *http.Request) {
		sess := req.Context().Value(server.SessionKey{}).(server.Session)
		body := req.Context().Value(server.BodyKey{}).(SetRequest)

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

	return middleware.Chain(handle,
		middleware.ValidateToken(ctx),
		middleware.ParseJSONBody[SetRequest](),
		middleware.Logging(ctx.Log),
	)
}

func setKey(s store.Store, user, key string, cipher, nonce []byte) error {
	return s.Set(user, key, store.CipherText{Nonce: nonce, Text: cipher})
}
