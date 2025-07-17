package errors

import "errors"

var InvalidRequestBody = errors.New("invalid request body")
var MissingUser = errors.New("\"username\" is required for logging in")
var MissingPassword = errors.New("\"password\" is required for logging in")
