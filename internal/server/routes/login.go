package routes

import (
	"encoding/json"
	"net/http"

	"github.com/jdpolicano/govault/internal/server"
)

type LoginSuccess struct {
	Token string `json:"sessionKey"`
}

type LoginRequest struct {
	Username string `json:"user"`
	Password string `json:"password"`
}

func GetLoginRoute(s *server.Context) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		encoder := json.NewEncoder(w)
		body, e := validateBody(req)
		if e != nil {
			encoder.Encode(loginFail(invalidRequestBody, 400))

		}
	}
}

func validateBody(req *http.Request) (LoginRequest, error) {
	var loginReq LoginRequest
	decoder := json.NewDecoder(req.Body)
	e := decoder.Decode(&loginReq)
	return loginReq, e
}

func loginFail(reason error, code int) LoginFailure {
	return LoginFailure{code, reason}
}
