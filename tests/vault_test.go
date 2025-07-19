package tests

import (
	"github.com/jdpolicano/govault/internal/vault"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key, err := vault.NewKey("password", 16)
	if err != nil {
		t.Fatalf("failed to create key: %v", err)
	}
	cipher, nonce, err := key.Encrypt("secret")
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}
	plain, err := key.Decrypt(nonce, string(cipher))
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}
	if string(plain) != "secret" {
		t.Errorf("expected decrypted text to match, got %q", plain)
	}
}

func TestGenerateRandBytes(t *testing.T) {
	b, err := vault.GenerateRandBytes(32)
	if err != nil {
		t.Fatalf("GenerateRandBytes returned error: %v", err)
	}
	if len(b) != 32 {
		t.Errorf("expected 32 bytes, got %d", len(b))
	}
}

func TestDeriveKeyWithSaltConsistency(t *testing.T) {
	salt, err := vault.GenerateRandBytes(16)
	if err != nil {
		t.Fatalf("failed to generate salt: %v", err)
	}
	k1, err := vault.NewKeyWithSalt("pass", salt)
	if err != nil {
		t.Fatalf("NewKeyWithSalt error: %v", err)
	}
	k2, err := vault.NewKeyWithSalt("pass", salt)
	if err != nil {
		t.Fatalf("NewKeyWithSalt error: %v", err)
	}
	if string(k1.Login) != string(k2.Login) || string(k1.AES) != string(k2.AES) {
		t.Errorf("keys derived with same salt differ")
	}
}
