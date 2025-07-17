package routes

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
func validateAuthRequest(req *http.Request) (AuthCredentials, error) {
	var reqData AuthCredentials
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&reqData)
	return reqData, err
}
