package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

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
func DeriveKeyFromText(text string, salt []byte) ([]byte, error) {
	// 600_000 was the recommended amount of iterations with HMAC-SHA-256
	// for password authentication. I think our usecase is more or less the same since
	// this key will be used to
	key, err := deriveKeyWithSalt(text, salt)
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

func deriveKeyWithSalt(text string, salt []byte) ([]byte, error) {
	// 600_000 was the recommended amount of iterations with HMAC-SHA-256
	// for password authentication. I think our usecase is more or less the same since
	// this key will be used to
	return pbkdf2.Key(sha256.New, text, salt, 600_000, 32)
}
