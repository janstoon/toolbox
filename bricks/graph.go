package bricks

import (
	"maps"
	"slices"
	"sync"

	"github.com/janstoon/toolbox/tricks"
)

// Graph is an adjacency map from keys to nodes. Weight is stored on each edge and is unused by
// BreadthFirstSearch, which treats every edge as unweighted.
type Graph[Key, Weight comparable, Value any] map[Key]*GraphNode[Key, Weight, Value]

// Add inserts a node for k with value v and returns that node. A second Add for the same key replaces the node.
func (g Graph[Key, Weight, Value]) Add(k Key, v Value) *GraphNode[Key, Weight, Value] {
	n := newGraphNode[Key, Weight, Value](v)
	g[k] = n

	return n
}

// Connect adds an undirected edge between root and each key in kk.
// Both endpoints must already exist. The edge weight is the zero value of Weight.
// Connect panics when either endpoint is missing.
func (g Graph[Key, Weight, Value]) Connect(root Key, kk ...Key) {
	var w Weight
	for _, key := range kk {
		g.connect(root, key, true, w)
	}
}

func (g Graph[Key, Weight, Value]) connect(k1, k2 Key, bidirect bool, weight Weight) {
	n1, ok := g[k1]
	if !ok {
		panic("connection origin not found")
	}

	n2, ok := g[k2]
	if !ok {
		panic("connection destination not found")
	}

	n1.connect(k2, weight)

	if bidirect {
		n2.connect(k1, weight)
	}
}

// BreadthFirstSearch returns one shortest path of keys from the start key to the goal key, including both ends.
// It returns nil when from is missing or no path exists.
func (g Graph[Key, Weight, Value]) BreadthFirstSearch(from, to Key) []Key {
	origin, ok := g[from]
	if !ok {
		return nil
	}

	var nn Queue[[]Key]
	nn.Enqueue(slices.Collect(tricks.Map(maps.Keys(origin.neighbors), func(k Key) []Key {
		return []Key{from, k}
	}))...)
	visited := make(map[Key]bool)

	for path := range nn.All() {
		k := path[len(path)-1]

		if _, ok := visited[k]; ok {
			continue
		}

		visited[k] = true

		if k == to {
			return path
		}

		n, ok := g[k]
		if !ok {
			panic("alien key")
		}

		nn.Enqueue(slices.Collect(tricks.Map(maps.Keys(n.neighbors), func(k Key) []Key {
			kk := make([]Key, len(path))
			copy(kk, path)

			return append(kk, k)
		}))...)
	}

	return nil
}

// GraphNode is one vertex. Value is the payload stored by Add.
type GraphNode[Key, Weight comparable, Value any] struct {
	neighbors map[Key]Weight
	// Value is the payload stored by Graph.Add.
	Value Value
}

func newGraphNode[Key, Weight comparable, Value any](v Value) *GraphNode[Key, Weight, Value] {
	return &GraphNode[Key, Weight, Value]{
		neighbors: make(map[Key]Weight),
		Value:     v,
	}
}

func (gn *GraphNode[Key, Weight, Value]) connect(key Key, weight Weight) {
	gn.neighbors[key] = weight
}

// TrieNode is one node of a prefix tree. Create the root with Trie.
type TrieNode[Key, KeyAtom comparable, Value any] struct {
	lock     sync.RWMutex
	children map[KeyAtom]*TrieNode[Key, KeyAtom, Value]
	atomizer func(Key) []KeyAtom

	full  bool
	value Value
}

// Trie returns an empty trie. atomizer splits a key into the atoms stored along the path,
// for example splitting a string into runes.
func Trie[Key, KeyAtom comparable, Value any](atomizer func(Key) []KeyAtom) *TrieNode[Key, KeyAtom, Value] {
	return &TrieNode[Key, KeyAtom, Value]{
		atomizer: atomizer,
	}
}

// Put stores val at key, replacing a value already stored at that key.
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

// Get returns the value stored at key, or nil when key was never put.
func (t *TrieNode[Key, _, Value]) Get(key Key) *Value {
	return t.get(t.atomizer(key))
}

func (t *TrieNode[_, KeyAtom, Value]) get(route []KeyAtom) *Value {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if len(route) == 0 {
		return t.val()
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

// BestMatch returns the value at the longest stored prefix of key.
// It returns nil when no prefix of key has a value.
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

// Delete clears the value stored at key. Missing keys are ignored.
func (t *TrieNode[Key, _, _]) Delete(key Key) {
	t.delete(t.atomizer(key))
}

func (t *TrieNode[_, KeyAtom, Value]) delete(route []KeyAtom) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if len(route) == 0 {
		var zero Value
		t.full = false
		t.value = zero

		return
	}

	if child, ok := t.children[route[0]]; ok {
		child.delete(route[1:])
	}
}
