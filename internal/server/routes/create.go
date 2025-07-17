package routes

import (
	"encoding/json"
	"fmt"
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

func GetCreateRoute(ctx *server.Context) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		body, e := validateBody(req)
		if e != nil {
			server.JSONResponse(w, server.NewInvalidBodyError())
			return
		}

		username, password := body.Username, body.Password
		user, salt, e := vault.CreateUserFromPW(username, password)
		if e != nil {
			ctx.Log.Println(e)
			server.JSONResponse(w, server.NewServerError(e))
			return
		}

		ctx.Log.Printf("successfully derived salt for user \"%s\"", username)

		if e := ctx.Store.AddUser(user); e != nil {
			ctx.Log.Println(e)
			switch e.(type) {
			case store.UserAlreadyExistsError:
				{
					server.JSONResponse(w, server.NewClientError(e))
				}
			default:
				{

					server.JSONResponse(w, server.NewServerError(e))
				}
			}
			return
		}

		ctx.Log.Printf("successfully added user \"%s\"", username)
		// create a session and return the token.
		sessId, idErr := server.GenerateSessionID()
		if idErr != nil {
			fmt.Println(idErr)
			server.JSONResponse(w, server.NewServerError(idErr))
			return
		}

		ctx.Log.Printf("starting session id=\"%s\"", sessId)
		aesKey, ciphErr := vault.CreateTextKeyFromPW("aes", password, salt)
		if ciphErr != nil {
			fmt.Println(idErr)
			server.JSONResponse(w, server.NewServerError(ciphErr))
			return
		}

		sess := server.NewSession(user.Name, aesKey, ctx.Config.DefaultTTL)
		ctx.Sessions.Set(sessId, sess)
		res := CreateSuccess{Token: sessId}
		server.JSONResponse(w, server.NewServerSuccess(res))
		ctx.Log.Println("response sent")
	}
}

func validateBody(req *http.Request) (CreateRequest, error) {
	var loginReq CreateRequest
	decoder := json.NewDecoder(req.Body)
	e := decoder.Decode(&loginReq)
	return loginReq, e
}
