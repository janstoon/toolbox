package tricks

import "errors"

var ErrDuplicateEntry = errors.New("duplicate entry")

func InsertIfNotExist[K comparable, V any](hm map[K]V, key K, value V) {
	_, exists := hm[key]
	if exists {
		panic(ErrDuplicateEntry)
	}

	hm[key] = value
}
