package server

import (
	"encoding/json"
	"net/http"
)

type TokenSuccess struct {
	Token string `json:"token"`
}

type AuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// validateRequest decodes and validates the request body.
func ValidateAuthRequest(req *http.Request) (AuthCredentials, error) {
	var reqData AuthCredentials
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqData)
	return reqData, err
}

// sendToken sends a successful response with the session token.
func SendToken(w http.ResponseWriter, token string) {
	res := TokenSuccess{Token: token}
	JSONResponse(w, NewServerSuccess(res))
}
