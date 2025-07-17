package routes

// import (
// 	"encoding/json"
// 	"errors"
// 	"net/http"

// 	"github.com/jdpolicano/govault/internal/server"
// )

// var invalidRequestBody = errors.New("invalid request body")
// var missingUser = errors.New("\"username\" is required for logging in")
// var missingPassword = errors.New("\"password\" is required for logging in")

// type LoginSuccess struct {
// 	Token string `json:"sessionKey"`
// }

// type LoginFailure struct {
// 	Code   int   `json:"code"`
// 	Reason error `json:"reason"`
// }

// type LoginRequest struct {
// 	Username string `json:"user"`
// 	Password string `json:"password"`
// }

// func GetLoginRoute(s *server.Context) func(w http.ResponseWriter, req *http.Request) {
// 	return func(w http.ResponseWriter, req *http.Request) {
// 		encoder := json.NewEncoder(w)
// 		body, e := validateBody(req)
// 		if e != nil {
// 			encoder.Encode(loginFail(invalidRequestBody, 400))

// 		}
// 	}
// }

// func validateBody(req *http.Request) (LoginRequest, error) {
// 	var loginReq LoginRequest
// 	decoder := json.NewDecoder(req.Body)
// 	e := decoder.Decode(&loginReq)
// 	return loginReq, e
// }

// func loginFail(reason error, code int) LoginFailure {
// 	return LoginFailure{code, reason}
// }
