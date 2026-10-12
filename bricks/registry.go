package bricks

import (
	"errors"
	"fmt"
	"sync"
)

// Registry is a goroutine-safe name-to-value map.
// Register rejects a name that is already present.
type Registry[T any] struct {
	l sync.RWMutex

	tt map[string]T
}

// NewRegistry returns an empty registry.
func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{
		tt: make(map[string]T),
	}
}

// Register stores entry under name.
// It returns an error that wraps ErrAlreadyExists when name is already registered.
func (reg *Registry[T]) Register(name string, entry T) error {
	reg.l.Lock()
	defer reg.l.Unlock()

	if _, ok := reg.tt[name]; ok {
		return errors.Join(fmt.Errorf("entry `%s` has been already registered", name), ErrAlreadyExists)
	}

	reg.tt[name] = entry

	return nil
}

// Get returns the entry stored under name.
// It returns ErrNotFound when name is absent.
func (reg *Registry[T]) Get(name string) (T, error) {
	reg.l.RLock()
	defer reg.l.RUnlock()

	t, ok := reg.tt[name]
	if !ok {
		return t, ErrNotFound
	}

	return t, nil
}

// All returns the entries currently stored. The caller must not mutate the map.
func (reg *Registry[T]) All() map[string]T {
	return reg.tt
}
