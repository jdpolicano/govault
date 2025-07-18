package errors

import (
	"errors"
	"fmt"
)

var InvalidRequestBody = errors.New("invalid request body")
var MissingUser = errors.New("\"username\" is required for logging in")
var MissingPassword = errors.New("\"password\" is required for logging in")
var IncorrectCredentials = errors.New("username or password incorrect")

func NewNoSuchUserError(u string) error {
	msg := fmt.Sprintf("no such user \"%s\"", u)
	return errors.New(msg)
}
