package store

import "encoding/base64"

type CipherText struct {
	Nonce string `json:"nonce"` // base64 encoded nonce
	Text  string `json:"text"`  // base64 encoded text
}

type User struct {
	Name  string `json:"name"`  // the name of the user
	Login string `json:"login"` // this is a namespaced pbkdf2 key derived from the pw for authenticaton purposes.
	Salt  string `json:"salt"`  // the salt for the login key and the aes key generation
}

func NewUser(name, login, salt string) User {
	return User{name, login, salt}
}

func (u User) SaltBytes() ([]byte, error) {
	l := base64.URLEncoding.DecodedLen(len(u.Salt))
	dest := make([]byte, l)
	_, err := base64.URLEncoding.Decode(dest, []byte(u.Salt))
	return dest, err
}

func NewUserFromBytes(name string, login, salt []byte) User {
	l := base64.URLEncoding.EncodeToString(login)
	s := base64.URLEncoding.EncodeToString(salt)
	return User{name, l, s}
}

type Store interface {
	GetUserInfo(string) (User, bool) // find info on a give user, if they exists
	AddUser(name, login, salt string) error
	HasUser(name string) bool
	Get(name, key string) (CipherText, bool)      // get a given key from the required key, nonce, and text.
	Set(name, key string, value CipherText) error // set a given value with a key
}
