package store

type Store[U any, V any] interface {
	GetUserInfo(string) (U, error)  // find info on a give user, if they exists
	Get(name, key string) (V, bool) // get a given key from the required key, nonce, and text.
	Set(name, key string) error     // set a given value with a key
}
