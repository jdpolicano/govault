package tests

import (
	"testing"
	"time"

	"github.com/jdpolicano/govault/internal/server"
)

func TestSessionExpiration(t *testing.T) {
	sess := server.NewSession("me", []byte("key"), time.Second)
	if sess.Expired() {
		t.Fatalf("session should not be expired immediately")
	}
	time.Sleep(1100 * time.Millisecond)
	if !sess.Expired() {
		t.Fatalf("session should be expired after ttl")
	}
}

func TestSessionMapOperations(t *testing.T) {
	sm := server.NewSessionMap()
	sm.Set("a", server.NewSession("u", []byte("k"), time.Minute))
	if _, ok := sm.Get("a"); !ok {
		t.Fatalf("expected to retrieve session")
	}
	sm.Delete("a")
	if _, ok := sm.Get("a"); ok {
		t.Fatalf("expected session to be deleted")
	}
}

func TestGenerateSessionID(t *testing.T) {
	id1, err := server.GenerateSessionID()
	if err != nil {
		t.Fatalf("error generating session id: %v", err)
	}
	id2, err := server.GenerateSessionID()
	if err != nil {
		t.Fatalf("error generating session id: %v", err)
	}
	if id1 == id2 {
		t.Errorf("expected unique session ids")
	}
	if len(id1) == 0 || len(id2) == 0 {
		t.Errorf("ids should not be empty")
	}
}
