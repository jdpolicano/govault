package routes

import (
	"errors"
	"net/http"

	"github.com/jdpolicano/govault/internal/server"
)

const missingUser = errors.New("\"username\" is required for logging in")
const missingPassword = errors.New("\"password\" is required for logging in")

type LoginRespone struct {
	token string `json:"sessionKey"`
}

type LoginRequest struct {
	username string `json:"user"`
	password string `json:"password"`
}

func GetLoginRoute(s *server.Server) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		req,
	}
}


func validateBody(body LoginRequest) error {
	return nil
}
