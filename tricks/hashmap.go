package tricks

import "errors"

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
