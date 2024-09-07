package bricks

import (
	"errors"
	"fmt"
	"sync"
)

type Registry[T any] struct {
	sync.RWMutex

	tt map[string]T
}

func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{
		tt: make(map[string]T),
	}
}

func (reg *Registry[T]) Register(name string, entry T) {
	reg.Lock()
	defer reg.Unlock()

	if _, ok := reg.tt[name]; ok {
		panic(errors.Join(fmt.Errorf("entry `%s` has been already registered", name), ErrAlreadyExists))
	}

	reg.tt[name] = entry
}

func (reg *Registry[T]) Get(name string) (T, error) {
	reg.RLock()
	defer reg.RUnlock()

	t, ok := reg.tt[name]
	if !ok {
		return t, ErrNotFound
	}

	return t, nil
}

func (reg *Registry[T]) All() map[string]T {
	return reg.tt
}
