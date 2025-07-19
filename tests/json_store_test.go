package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jdpolicano/govault/internal/store"
)

func TestJSONStoreUserLifecycle(t *testing.T) {
	dir := t.TempDir()
	js := store.NewJSONStore(dir)

	if js.HasUser("bob") {
		t.Fatalf("expected store to have no users")
	}

	err := js.AddUser("bob", []byte("login"), []byte("salt"))
	if err != nil {
		t.Fatalf("AddUser returned error: %v", err)
	}

	if !js.HasUser("bob") {
		t.Errorf("expected HasUser to return true")
	}

	u, ok := js.GetUserInfo("bob")
	if !ok || u.Name != "bob" {
		t.Fatalf("GetUserInfo returned wrong user")
	}

	// check file exists
	path := filepath.Join(dir, "bob", "secrets.json")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected user file to exist: %v", err)
	}

	// test Set and Get
	err = js.Set("bob", "key", store.CipherText{Nonce: []byte("n"), Text: []byte("c")})
	if err != nil {
		t.Fatalf("Set returned error: %v", err)
	}
	ct, ok := js.Get("bob", "key")
	if !ok || !ct.Equal(store.CipherText{Nonce: []byte("n"), Text: []byte("c")}) {
		t.Fatalf("Get returned wrong value")
	}
}
