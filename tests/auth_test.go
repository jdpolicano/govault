package tests

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/jdpolicano/govault/internal/server"
)

func TestValidateAuthRequest(t *testing.T) {
	body := bytes.NewBufferString(`{"username":"u","password":"p"}`)
	req := httptest.NewRequest("POST", "/", body)
	cred, err := server.ValidateAuthRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cred.Username != "u" || cred.Password != "p" {
		t.Errorf("unexpected credentials: %+v", cred)
	}
}
