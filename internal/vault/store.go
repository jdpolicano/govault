package vault

import (
	"encoding/base64"

	"github.com/jdpolicano/govault/internal/store"
)

func CreateTextKeyFromPW(namespace, pw string, salt []byte) (string, error) {
	k := namespace + pw
	key, err := DeriveKeyFromText(k, salt)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(key), err
}

func CreateUserFromPW(username, pw string) (store.User, []byte, error) {
	var user store.User

	salt, err := GenerateRandBytes(16)
	if err != nil {
		return user, nil, err
	}

	namespace, err := CreateTextKeyFromPW("login", pw, salt)
	if err != nil {
		return user, nil, err
	}

	key, err := DeriveKeyFromText(namespace, salt)
	if err != nil {
		return user, nil, err
	}

	return store.NewUserFromBytes(username, key, salt), salt, nil
}
