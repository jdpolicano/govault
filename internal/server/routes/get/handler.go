package get

import (
	"net/http"

	"github.com/jdpolicano/govault/internal/server"
	e "github.com/jdpolicano/govault/internal/server/errors"
	"github.com/jdpolicano/govault/internal/server/middleware"
	"github.com/jdpolicano/govault/internal/vault"
)

type GetRequest struct {
	Key string `json:"key"`
}

func Handler(ctx *server.Context) http.HandlerFunc {
	handle := func(w http.ResponseWriter, req *http.Request) {
		sess := req.Context().Value(server.SessionKey{}).(server.Session)
		body := req.Context().Value(server.BodyKey{}).(GetRequest)

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

	return middleware.Chain(handle,
		middleware.ValidateToken(ctx),
		middleware.ParseJSONBody[GetRequest](),
		middleware.Logging(ctx.Log),
	)
}
