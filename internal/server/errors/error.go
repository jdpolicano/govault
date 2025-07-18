package errors

import (
	"errors"
	"fmt"
)

var InvalidRequestBody = errors.New("invalid request body")
var MissingUser = errors.New("\"username\" is required for logging in")
var MissingPassword = errors.New("\"password\" is required for logging in")
var IncorrectCredentials = errors.New("username or password incorrect")
var MissingAuthorizationHeader = errors.New("Authorization Failed")
var MalformedAuthorizationHeader = errors.New("Authorization Header Malformed")
var AuthorizationExpired = errors.New("Authorization Header Malformed")
var UnexpectedServerError = errors.New("Unexpected Serverside Error")

func NewNoSuchUserError(u string) error {
	msg := fmt.Sprintf("no such user \"%s\"", u)
	return errors.New(msg)
}
