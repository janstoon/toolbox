package bricks

import (
	"iter"
	"sync"
)

type Queue[T any] struct {
	lock sync.RWMutex

	tt []T
}

func (q *Queue[T]) Enqueue(tt ...T) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.tt = append(q.tt, tt...)
}

func (q *Queue[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		itr := q.iter()
		for v, ok := itr.next(); ok; v, ok = itr.next() {
			if !yield(v) {
				return
			}
		}
	}
}

func (q *Queue[T]) iter() *walker[T] {
	return &walker[T]{
		get: func(i int) (T, bool) {
			q.lock.RLock()
			defer q.lock.RUnlock()

			if i < len(q.tt) && i > -1 {
				return q.tt[i], true
			}

			var t T

			return t, false
		},
	}
}

type walker[T any] struct {
	counter int
	get     func(i int) (T, bool)
}

func (w *walker[T]) next() (T, bool) {
	t, ok := w.get(w.counter)
	if ok {
		w.counter++
	}

	return t, ok
}
