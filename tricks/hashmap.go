package tricks

import "errors"

// ErrDuplicateEntry is the panic value of InsertIfNotExist when key is already present.
var ErrDuplicateEntry = errors.New("duplicate entry")

// InsertIfNotExist puts value in hm[key] if hm[key] is not already full.
// If it's already filled it panics with ErrDuplicateEntry
func InsertIfNotExist[K comparable, V any](hm map[K]V, key K, value V) {
	_, exists := hm[key]
	if exists {
		panic(ErrDuplicateEntry)
	}

	hm[key] = value
}
