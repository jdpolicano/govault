package store

import "fmt"

type UserAlreadyExistsError struct {
	name string
}

func NewAlreadyExistsError(name string) UserAlreadyExistsError {
	return UserAlreadyExistsError{name}
}

func (e UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("err user %s already exists", e.name)
}
