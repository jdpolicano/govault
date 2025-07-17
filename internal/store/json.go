package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type CipherText struct {
	Nonce string `json:"nonce"` // base64 encoded nonce
	Text  string `json:"text"`  // base64 encoded text
}

type User struct {
	Name  string `json:"name"`  // the name of the user
	Login string `json:"login"` // the namespaced key for login comparison
	Salt  string `json:"salt"`  // the salt for the login key and the aes key generation
}

type StoreRecord struct {
	user    User
	secrets map[string]CipherText
}

// in memory and file backed json store.
type JSONStore struct {
	sync.RWMutex
	vaultPath string                 // the path to the store's location
	data      map[string]StoreRecord // the in memory store, backed by the disk map
	disk      map[string]*os.File    // these are the actual file descriptors that we have opened so we can quickly update a write over the given file
}

func NewJSONStore(path string) *JSONStore {
	return &JSONStore{
		vaultPath: path,
		data:      make(map[string]StoreRecord, 1024),
		disk:      make(map[string]*os.File, 1024),
	}
}

func (js *JSONStore) GetUserInfo(name string) (User, bool) {
	js.RLock()
	defer js.RUnlock()
	var none User
	record, exists := js.data[name]
	if !exists {
		return none, false
	}
	return record.user, true
}

func (js *JSONStore) Get(name, key string, value CipherText) (CipherText, bool) {
	js.RLock()
	defer js.RUnlock()
	var none CipherText
	record, recExists := js.data[name]
	if !recExists {
		return none, false
	}
	ciphertext, ciphExists := record.secrets[key]
	if !ciphExists {
		return none, false
	}
	return ciphertext, true
}

func (js *JSONStore) Set(name, key string, value CipherText) error {
	js.Lock()
	defer js.Unlock()
	record, userExists := js.data[name]
	if !userExists {
		return fmt.Errorf("err user %s does not exist", name)
	}
	currCipher, cipherExists := record.secrets[key]
	if cipherExists && currCipher == value {
		return nil
	}
	record.secrets[key] = value
	userP := js.getUserPath(record.user.Name)
	return recordOnDisk(userP, record)
}

func (js *JSONStore) getUserPath(user string) string {
	return filepath.Join(js.vaultPath, user, "secrets.json")
}

func recordOnDisk(path string, record StoreRecord) error {
	keys := make([]CipherText, 0, len(record.secrets))
	for _, cipherText := range record.secrets {
		keys = append(keys, cipherText)
	}
	bytes, err := json.Marshal(keys)
	if err != nil {
		return err
	}
	// ensure the parent directories exist
	dirErr := os.MkdirAll(filepath.Dir(path), 0700)
	if dirErr != nil {
		return dirErr
	}
	// write the file
	return os.WriteFile(path, bytes, 0700)
}
