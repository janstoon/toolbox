package bricks

import (
	"errors"
	"fmt"
	"sync"
)

type Registry[T any] struct {
	l sync.RWMutex

	tt map[string]T
}

func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{
		tt: make(map[string]T),
	}
}

func (reg *Registry[T]) Register(name string, entry T) error {
	reg.l.Lock()
	defer reg.l.Unlock()

	if _, ok := reg.tt[name]; ok {
		return errors.Join(fmt.Errorf("entry `%s` has been already registered", name), ErrAlreadyExists)
	}

	reg.tt[name] = entry

	return nil
}

func (reg *Registry[T]) Get(name string) (T, error) {
	reg.l.RLock()
	defer reg.l.RUnlock()

	t, ok := reg.tt[name]
	if !ok {
		return t, ErrNotFound
	}

	return t, nil
}

func (reg *Registry[T]) All() map[string]T {
	return reg.tt
}
