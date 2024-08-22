package tricks

type Iterator[T comparable] interface {
	Next() (T, bool)
}

type SliceIterator[T comparable] struct {
	slice []T
	index int
}

func NewSliceIterator[T comparable](s []T) *SliceIterator[T] {
	return &SliceIterator[T]{slice: s}
}

func (si *SliceIterator[T]) Next() (T, bool) {
	if si.index >= len(si.slice) {
		var dst T
		return dst, false
	}

	v := si.slice[si.index]
	si.index++

	return v, true
}

// Iter returns an iterator extension.
// The extension provides additional functionality for the iterator.
// The extension is not thread-safe.
// The extension is not reusable. If you need to iterate again, call Clone.
func (si *SliceIterator[T]) Iter() iteratorExtension[T] {
	si.index = 0
	return iteratorExtension[T]{i: si}
}

type iteratorExtension[T comparable] struct {
	i Iterator[T]
}

func (im iteratorExtension[T]) Next() (T, bool) {
	return im.i.Next()
}

// All returns true if all elements in the iterator satisfy the predicate.
// If the iterator is empty, it returns true.
func (im iteratorExtension[T]) All(f func(T) bool) bool {
	for {
		v, ok := im.i.Next()
		if !ok {
			break
		}

		if !f(v) {
			return false
		}
	}

	return true
}

// Any returns true if any element in the iterator satisfies the predicate.
// If the iterator is empty, it returns false.
func (im iteratorExtension[T]) Any(f func(T) bool) bool {
	for {
		v, ok := im.i.Next()
		if !ok {
			break
		}

		if f(v) {
			return true
		}
	}

	return false
}

// Contains returns true if the iterator contains the value.
func (im iteratorExtension[T]) Contains(v T) bool {
	for {
		i, ok := im.i.Next()
		if !ok {
			break
		}

		if i == v {
			return true
		}
	}

	return false
}

// ContainsBy returns true if the iterator contains the value.
// The function f is used to compare the value.
func (im iteratorExtension[T]) ContainsBy(f func(T) bool) bool {
	for {
		i, ok := im.i.Next()
		if !ok {
			break
		}

		if f(i) {
			return true
		}
	}

	return false
}

// Find returns the first element that satisfies the predicate.
// If no element satisfies the predicate, it returns the zero value of T and false.
func (im iteratorExtension[T]) Find(f func(T) bool) (T, bool) {
	for {
		v, ok := im.i.Next()
		if !ok {
			break
		}

		if f(v) {
			return v, true
		}
	}

	var z T
	return z, false
}

// FindIndex returns the index of the first element that satisfies the predicate.
// If no element satisfies the predicate, it returns -1.
// The index is zero-based.
func (im iteratorExtension[T]) FindIndex(f func(T) bool) int {
	i := 0

	for {
		v, ok := im.i.Next()
		if !ok {
			break
		}

		if f(v) {
			return i
		}

		i++
	}

	return -1
}

// FindLastIndex returns the index of the last element that satisfies the predicate.
func (im iteratorExtension[T]) FindLastIndex(f func(T) bool) int {
	i := -1
	index := 0

	for {
		v, ok := im.i.Next()
		if !ok {
			break
		}

		if f(v) {
			i = index
		}

		index++
	}

	return i
}

// Reduce reduces the iterator to a single value.
// The function f is used to combine the elements.
// The first element is used as the initial value.
// If the iterator is empty, it returns the zero value of T.
func (im iteratorExtension[T]) Reduce(f func(T, T) T) T {
	var r T
	first := true

	for {
		v, ok := im.i.Next()
		if !ok {
			break
		}

		if first {
			r = v
			first = false
		} else {
			r = f(r, v)
		}
	}

	return r
}

// Filter returns a new iterator with elements that satisfy the predicate.
// The order of the elements is preserved.
// If no element satisfies the predicate, it returns an empty iterator.
func (im iteratorExtension[T]) Filter(f func(T) bool) iteratorExtension[T] {
	var d []T

	for {
		v, ok := im.i.Next()
		if !ok {
			break
		}

		if f(v) {
			d = append(d, v)
		}
	}

	return NewSliceIterator(d).Iter()
}

// Map returns a new iterator with elements that are transformed by the function.
// Because of the poor go generics support, the return type is any.
// Please wrap the iterator with IteratorCast to cast the elements to the desired type.
func (im iteratorExtension[T]) Map(f func(T) any) iteratorExtension[any] {
	var d []any

	for {
		v, ok := im.i.Next()
		if !ok {
			break
		}

		d = append(d, f(v))
	}

	return NewSliceIterator(d).Iter()
}

// FlatMap returns a new iterator with elements that are transformed by the function.
func (im iteratorExtension[T]) FlatMap(f func(T) []T) iteratorExtension[T] {
	var d []T

	for {
		v, ok := im.i.Next()
		if !ok {
			break
		}

		d = append(d, f(v)...)
	}

	return NewSliceIterator(d).Iter()
}

// Distinct returns a new iterator with distinct elements.
func (im iteratorExtension[T]) Distinct() iteratorExtension[T] {
	m := make(map[T]struct{})
	var d []T

	for {
		v, ok := im.i.Next()
		if !ok {
			break
		}

		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			d = append(d, v)
		}
	}

	return NewSliceIterator(d).Iter()
}

// DistinctBy returns a new iterator with distinct elements.
// The function f is used to determine the key of the element.
func (im iteratorExtension[T]) DistinctBy(f func(T) T) iteratorExtension[T] {
	m := make(map[T]struct{})
	var d []T

	for {
		v, ok := im.i.Next()
		if !ok {
			break
		}

		key := f(v)
		if _, ok := m[key]; !ok {
			m[key] = struct{}{}
			d = append(d, v)
		}
	}

	return NewSliceIterator(d).Iter()
}

// First returns the first element in the iterator.
// If the iterator is empty, it returns the zero value of T and false.
func (im iteratorExtension[T]) First() (T, bool) {
	return im.i.Next()
}

// Last returns the last element in the iterator.
// If the iterator is empty, it returns the zero value of T and false.
func (im iteratorExtension[T]) Last() (T, bool) {
	var v T
	var ok bool

	for {
		var oki bool
		var vi T
		vi, oki = im.i.Next()
		if !oki {
			break
		}

		v = vi
		ok = true
	}

	return v, ok
}

// ElementAt returns the element at the specified index.
// If the index is out of range, it returns the zero value of T and false.
func (im iteratorExtension[T]) ElementAt(index int) (T, bool) {
	i := 0

	for {
		v, ok := im.i.Next()
		if !ok {
			break
		}

		if i == index {
			return v, true
		}

		i++
	}

	var dst T
	return dst, false
}

// ElementAtOrDefault returns the element at the specified index.
// If the index is out of range, it returns the default value.
func (im iteratorExtension[T]) ElementAtOrDefault(index int, defaultValue T) T {
	v, ok := im.ElementAt(index)
	if !ok {
		return defaultValue
	}

	return v
}

// ElementAtOrElse returns the element at the specified index.
// If the index is out of range, it returns the value from the function.
func (im iteratorExtension[T]) ElementAtOrElse(index int, f func() T) T {
	v, ok := im.ElementAt(index)
	if !ok {
		return f()
	}

	return v
}

// Collect returns a slice with all elements in the iterator.
func (im iteratorExtension[T]) Collect() []T {
	var d []T = make([]T, 0)

	for {
		v, ok := im.i.Next()
		if !ok {
			break
		}

		d = append(d, v)
	}

	return d
}

// Count returns the number of elements in the iterator.
func (im iteratorExtension[T]) Count() int {
	c := 0

	for {
		_, ok := im.i.Next()
		if !ok {
			break
		}

		c++
	}

	return c
}

// ForEach calls the function for each element in the iterator.
func (im iteratorExtension[T]) ForEach(f func(T)) {
	for {
		v, ok := im.i.Next()
		if !ok {
			break
		}

		f(v)
	}
}

// IndexOf returns the index of the element.
// If the element is not found, it returns -1.
func (im iteratorExtension[T]) IndexOf(v T) int {
	i := 0

	for {
		e, ok := im.i.Next()
		if !ok {
			break
		}

		if e == v {
			return i
		}

		i++
	}

	return -1
}

// LastIndexOf returns the last index of the element.
// If the element is not found, it returns -1.
func (im iteratorExtension[T]) LastIndexOf(v T) int {
	i := -1
	index := 0

	for {
		e, ok := im.i.Next()
		if !ok {
			break
		}

		if e == v {
			i = index
		}

		index++
	}

	return i
}

// Skip returns a new iterator with the first n elements skipped.
func (im iteratorExtension[T]) Skip(n int) iteratorExtension[T] {
	for i := 0; i < n; i++ {
		_, ok := im.i.Next()
		if !ok {
			break
		}
	}

	return im
}

// Take returns a new iterator with the first n elements taken.
func (im iteratorExtension[T]) Take(n int) iteratorExtension[T] {
	var d []T

	for i := 0; i < n; i++ {
		v, ok := im.i.Next()
		if !ok {
			break
		}

		d = append(d, v)
	}

	return NewSliceIterator(d).Iter()
}

// TakeWhile returns a new iterator with elements taken while the predicate is true.
func (im iteratorExtension[T]) TakeWhile(f func(T) bool) iteratorExtension[T] {
	var d []T = make([]T, 0)

	var indx = im.Clone().FindLastIndex(f)
	if indx == -1 {
		return NewSliceIterator(d).Iter()
	}

	for i := 0; i <= indx; i++ {
		v, ok := im.i.Next()
		if !ok {
			break
		}

		d = append(d, v)
	}

	return NewSliceIterator(d).Iter()
}

// SkipWhile returns a new iterator with elements skipped while the predicate is true.
func (im iteratorExtension[T]) SkipWhile(f func(T) bool) iteratorExtension[T] {
	var indx = im.Clone().FindLastIndex(f)
	if indx == -1 {
		return NewSliceIterator([]T{}).Iter()
	}

	for i := 0; i <= indx; i++ {
		im.i.Next()
	}

	return im
}

// Chain returns a new iterator that chains the current iterator with another iterator.
func (im iteratorExtension[T]) Chain(i Iterator[T]) iteratorExtension[T] {
	return iteratorExtension[T]{i: NewChainIterator(im.i, i)}
}

// Extend returns a new iterator that extends the current iterator with a slice.
func (im iteratorExtension[T]) Extend(i []T) iteratorExtension[T] {
	return iteratorExtension[T]{i: NewChainIterator(im.i, NewSliceIterator(i))}
}

// Clone returns a new iterator that is a clone of the current iterator.
func (im *iteratorExtension[T]) Clone() iteratorExtension[T] {
	var cur = im.Collect()
	im.i = NewSliceIterator(cur)

	return NewSliceIterator(cur).Iter()
}

// Go 1.22 experimental and 1.23 language support for range over func
// see: https://go.dev/wiki/RangefuncExperiment
// and: https://go.dev/doc/go1.23#language
func (im *iteratorExtension[T]) Range() func(yield func(T) bool) {
	return func(yield func(T) bool) {
		for {
			v, ok := im.i.Next()
			if !ok {
				break
			}

			if !yield(v) {
				break
			}
		}
	}
}

// Go 1.22 experimental and 1.23 language support for range over func
// see: https://go.dev/wiki/RangefuncExperiment
// and: https://go.dev/doc/go1.23#language
func (im *iteratorExtension[T]) Enumerate() func(yield func(int, T) bool) {
	return func(yield func(int, T) bool) {
		i := 0
		for {
			v, ok := im.i.Next()
			if !ok {
				break
			}

			if !yield(i, v) {
				break
			}

			i++
		}
	}
}

type ChainIterator[T comparable] struct {
	i2      []Iterator[T]
	current int
}

func NewChainIterator[T comparable](is ...Iterator[T]) *ChainIterator[T] {
	return &ChainIterator[T]{i2: is}
}

func (ci *ChainIterator[T]) Next() (T, bool) {
	for ci.current < len(ci.i2) {
		v, ok := ci.i2[ci.current].Next()
		if ok {
			return v, true
		}

		ci.current++
	}

	var dst T
	return dst, false
}

// IteratorCast casts the elements of the iterator to the specified type.
// It is useful in cases where the type is lost in the transformation.
func IteratorCast[T comparable](i Iterator[any]) iteratorExtension[T] {
	return (&castIterator[T]{i: i}).Iter()
}

type castIterator[T comparable] struct {
	i Iterator[any]
}

func (ci *castIterator[T]) Next() (T, bool) {
	v, ok := ci.i.Next()
	if !ok {
		var dst T
		return dst, false
	}

	return v.(T), true
}

func (ci *castIterator[T]) Iter() iteratorExtension[T] {
	return iteratorExtension[T]{i: ci}
}
