package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type JSONRecord struct {
	User    User                  `json:"user"`
	Secrets map[string]CipherText `json:"secrets"`
}

func NewJSONRecord(user User) JSONRecord {
	return JSONRecord{user, make(map[string]CipherText, 256)}
}

// in memory and file backed json store.
type JSONStore struct {
	sync.RWMutex
	vaultPath string                // the path to the store's location
	data      map[string]JSONRecord // the in memory store, backed by a json file
}

// todo: since we have the vault path, we should setup the store based on whatever was there before!
// todo: we should have some kind eviction policy since we can't assume we'll hold all of these secrets in memory
func NewJSONStore(path string) *JSONStore {
	return &JSONStore{
		vaultPath: path,
		data:      make(map[string]JSONRecord, 1024),
	}
}

func (js *JSONStore) AddUser(name string, login, salt []byte) error {
	js.Lock()
	defer js.Unlock()
	record, exists := js.data[name]
	if exists {
		return NewAlreadyExistsError(name)
	}
	record = NewJSONRecord(User{name, login, salt})
	userP := js.getUserPath(record.User.Name)
	if e := recordOnDisk(userP, record); e != nil {
		return e
	}
	js.data[name] = record
	return nil
}

func (js *JSONStore) HasUser(name string) bool {
	js.RLock()
	defer js.RUnlock()
	_, exists := js.data[name]
	if !exists {
		return false
	}
	return true
}

func (js *JSONStore) GetUserInfo(name string) (User, bool) {
	js.RLock()
	defer js.RUnlock()
	var none User
	record, exists := js.data[name]
	if !exists {
		return none, false
	}
	return record.User, true
}

func (js *JSONStore) Get(name, key string) (CipherText, bool) {
	js.RLock()
	defer js.RUnlock()
	var none CipherText
	record, recExists := js.data[name]
	if !recExists {
		return none, false
	}
	ciphertext, ciphExists := record.Secrets[key]
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

	original, cipherExists := record.Secrets[key]
	if cipherExists && original.Equal(value) {
		return nil
	}

	record.Secrets[key] = value
	userP := js.getUserPath(record.User.Name)
	if e := recordOnDisk(userP, record); e != nil {
		record.Secrets[key] = original
		return e
	}
	return nil
}

func (js *JSONStore) getUserPath(user string) string {
	return filepath.Join(js.vaultPath, user, "secrets.json")
}

func recordOnDisk(path string, record JSONRecord) error {
	bytes, err := json.Marshal(record)
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
