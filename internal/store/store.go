package store

import "bytes"

type CipherText struct {
	Nonce []byte `json:"nonce"` // base64 encoded nonce
	Text  []byte `json:"text"`  // base64 encoded text
}

func (c CipherText) Equal(other CipherText) bool {
	return bytes.Equal(c.Nonce, other.Nonce) && bytes.Equal(c.Text, other.Text)
}

type User struct {
	Name  string `json:"name"`  // the name of the user
	Login []byte `json:"login"` // this is a namespaced pbkdf2 key derived from the pw for authenticaton purposes.
	Salt  []byte `json:"salt"`  // the salt for the login key and the aes key generation
}

func NewUser(name string, login, salt []byte) User {
	return User{name, login, salt}
}

type Store interface {
	GetUserInfo(string) (User, bool) // find info on a give user, if they exists
	AddUser(name string, login, salt []byte) error
	HasUser(name string) bool
	Get(name, key string) (CipherText, bool)      // get a given key from the required key, nonce, and text.
	Set(name, key string, value CipherText) error // set a given value with a key
}
