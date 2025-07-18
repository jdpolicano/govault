package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

type Key struct {
	Login []byte
	AES   []byte
	Salt  []byte
}

func NewKey(text string, saltSize int) (*Key, error) {
	salt, err := GenerateRandBytes(saltSize)
	if err != nil {
		return nil, err
	}

	login, err := pbkdf2Key("login:"+text, salt)
	if err != nil {
		return nil, err
	}

	aes, err := pbkdf2Key("aes:"+text, salt)
	if err != nil {
		return nil, err
	}

	return &Key{login, aes, salt}, nil
}

func KeyFromSaltString(text string, saltStr string) (*Key, error) {
	bytes, err := base64.RawStdEncoding.DecodeString(saltStr)
	fmt.Println("raw bytes", bytes)
	if err != nil {
		return nil, err
	}
	login, err := pbkdf2Key("login:"+text, bytes)
	if err != nil {
		return nil, err
	}

	fmt.Println("raw login bytes", bytes)
	aes, err := pbkdf2Key("aes:"+text, bytes)
	if err != nil {
		return nil, err
	}

	return &Key{login, aes, bytes}, nil
}

func (k *Key) Encrypt(plaintext string) ([]byte, []byte, error) {
	return Encrypt(k.AES, plaintext)
}

func (k *Key) Decrypt(nonce []byte, cipherText string) ([]byte, error) {
	return Decrypt(nonce, k.AES, []byte(cipherText))
}

func (k *Key) Base64LoginKey() string {
	return base64.RawStdEncoding.EncodeToString(k.Login)
}

func (k *Key) Base64Salt() string {
	return base64.RawStdEncoding.EncodeToString(k.Salt)
}

// Generate n random bytes
func GenerateRandBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to read random bytes: %w", err)
	}
	return b, nil
}

// takes a plaintext string and returns a base64 encoded key and a base64 encoded salt
// errors if its unable to generate a random salt or the pbkdf2 func fails.
func DeriveKey(text string, saltSize int) ([]byte, []byte, error) {
	// 600_000 was the recommended amount of iterations with HMAC-SHA-256
	// for password authentication. I think our usecase is more or less the same since
	// this key will be used to
	salt, err := GenerateRandBytes(saltSize)
	if err != nil {
		return nil, nil, err
	}
	key, err := pbkdf2Key(text, salt)
	return key, salt, err
}

// takes a plaintext string and returns a base64 encoded key and a base64 encoded salt
// errors if its unable to generate a random salt or the pbkdf2 func fails.
func DeriveKeyWithSalt(text string, salt []byte) ([]byte, error) {
	// 600_000 was the recommended amount of iterations with HMAC-SHA-256
	// for password authentication. I think our usecase is more or less the same since
	// this key will be used to
	key, err := pbkdf2Key(text, salt)
	return key, err
}

func Encrypt(key []byte, text string) ([]byte, []byte, error) {
	aesgcm, err := createGCM(key)
	if err != nil {
		return nil, nil, err
	}
	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	// todo: we should keep track of the nonce's we use with each unique key
	nonce, err := GenerateRandBytes(aesgcm.NonceSize())
	if err != nil {
		return nil, nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(text), nil)
	return ciphertext, nonce, nil
}

func Decrypt(nonce, key, ciphertext []byte) ([]byte, error) {
	aesgcm, err := createGCM(key)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func createGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aesgcm, err
}

func pbkdf2Key(text string, salt []byte) ([]byte, error) {
	// 600_000 was the recommended amount of iterations with HMAC-SHA-256
	// for password authentication. I think our usecase is more or less the same since
	// this key will be used to
	return pbkdf2.Key(sha256.New, text, salt, 600_000, 32)
}
