package bricks

import "sync"

type TrieNode[Key, KeyAtom comparable, Value any] struct {
	lock     sync.RWMutex
	children map[KeyAtom]*TrieNode[Key, KeyAtom, Value]
	atomizer func(Key) []KeyAtom

	full  bool
	value Value
}

func Trie[Key, KeyAtom comparable, Value any](atomizer func(Key) []KeyAtom) *TrieNode[Key, KeyAtom, Value] {
	return &TrieNode[Key, KeyAtom, Value]{
		atomizer: atomizer,
	}
}

func (t *TrieNode[Key, KeyAtom, Value]) Put(key Key, val Value) {
	t.put(t.atomizer(key), val)
}

func (t *TrieNode[Key, KeyAtom, Value]) put(route []KeyAtom, val Value) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if len(route) == 0 {
		t.full = true
		t.value = val

		return
	}

	if t.children == nil {
		t.children = make(map[KeyAtom]*TrieNode[Key, KeyAtom, Value])
	}

	if _, ok := t.children[route[0]]; !ok {
		t.children[route[0]] = Trie[Key, KeyAtom, Value](t.atomizer)
	}

	t.children[route[0]].put(route[1:], val)
}

func (t *TrieNode[Key, _, Value]) Get(key Key) *Value {
	return t.get(t.atomizer(key))
}

func (t *TrieNode[_, KeyAtom, Value]) get(route []KeyAtom) *Value {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if len(route) == 0 {
		return &t.value
	}

	if child, ok := t.children[route[0]]; ok {
		return child.get(route[1:])
	}

	return nil
}

func (t *TrieNode[_, KeyAtom, Value]) val() *Value {
	if t.full {
		return &t.value
	}

	return nil
}

func (t *TrieNode[Key, _, Value]) BestMatch(key Key) *Value {
	return t.bestMatch(t.atomizer(key))
}

func (t *TrieNode[_, KeyAtom, Value]) bestMatch(route []KeyAtom) *Value {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if len(route) == 0 {
		return t.val()
	}

	if child, ok := t.children[route[0]]; ok {
		if v := child.bestMatch(route[1:]); v != nil {
			return v
		}
	}

	return t.val()
}

func (t *TrieNode[Key, _, _]) Delete(key Key) {
	panic("not implemented")
}
