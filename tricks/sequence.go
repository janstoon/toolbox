package tricks

import (
	"iter"
	"slices"
)

// IndexUnavailable is used in methods which are ought to return index of something, but they failed finding it and
// want to express that it's unavailable,
const IndexUnavailable = -1

// SliceClone creates a copy of slice(s) which can be modified independently without any effect on original slice.
// It's equivalent to slices.Clone
func SliceClone[Slice ~[]E, E any](s Slice) Slice {
	cloned := make(Slice, len(s))
	copy(cloned, s)

	return cloned
}

// All returns true if all elements in the seq satisfy the predicate using matcher.
// If seq is empty, it returns true.
func All[E any](seq iter.Seq[E], matcher Matcher[E]) bool {
	for v := range seq {
		if !matcher.Match(v) {
			return false
		}
	}

	return true
}

// SliceAll performs All on slice
func SliceAll[Slice ~[]E, E any](s Slice, matcher Matcher[E]) bool {
	return All(slices.Values(s), matcher)
}

// Any returns true if any element in the seq satisfies the predicate using matcher.
// If the seq is empty, it returns false. It's equivalent to slices.ContainsFunc
func Any[E any](seq iter.Seq[E], matcher Matcher[E]) bool {
	for v := range seq {
		if matcher.Match(v) {
			return true
		}
	}

	return false
}

// SliceAny performs Any on slice
func SliceAny[Slice ~[]E, E any](s Slice, matcher Matcher[E]) bool {
	return Any(slices.Values(s), matcher)
}

// Find returns the first element in seq that satisfies the predicate using matcher.
// If no element was found, it returns nil.
func Find[E any](seq iter.Seq[E], matcher Matcher[E]) (E, bool) {
	for e := range seq {
		if matcher.Match(e) {
			return e, true
		}
	}

	var e E

	return e, false
}

// SliceFind performs Find on slice
func SliceFind[Slice ~[]E, E any](s Slice, matcher Matcher[E]) (E, bool) {
	return Find(slices.Values(s), matcher)
}

// Index returns the index of the first element that satisfies the predicate using matcher and true.
// If no element was found, it returns zero value of K and false.
func Index[K, V any](seq iter.Seq2[K, V], matcher Matcher[V]) (K, bool) {
	for k, v := range seq {
		if matcher.Match(v) {
			return k, true
		}
	}

	var i K

	return i, false
}

// SliceIndex performs Index on slice. It returns IndexUnavailable if no element was found.
// It's Equivalent to slices.IndexFunc.
func SliceIndex[Slice ~[]E, E any](s Slice, matcher Matcher[E]) int {
	i, ok := Index(slices.All(s), matcher)
	if !ok {
		return IndexUnavailable
	}

	return i
}

// IndexOf returns the index of the element and true.
// If no element was found, it returns zero value of K and false.
func IndexOf[K any, V comparable](seq iter.Seq2[K, V], expected V) (K, bool) {
	return Index(seq, MatchEqual(expected))
}

// SliceIndexOf searches slice for first element which is equal to expected and returns index of it
// or IndexUnavailable if not found.
func SliceIndexOf[Slice ~[]E, E comparable](expected E, s Slice) int {
	return SliceIndex(s, MatchEqual(expected))
}

// LastIndex returns the index of the last element that satisfies the predicate using matcher and true.
// If no element was found, it returns zero value of K and false.
func LastIndex[K, V any](seq iter.Seq2[K, V], matcher Matcher[V]) (K, bool) {
	var i K
	ok := false

	for k, v := range seq {
		if matcher.Match(v) {
			i, ok = k, true
		}
	}

	return i, ok
}

// SliceLastIndex performs LastIndex on slice. It returns IndexUnavailable if no element was found.
// It's Equivalent to slices.IndexFunc.
func SliceLastIndex[Slice ~[]E, E any](s Slice, matcher Matcher[E]) int {
	i, ok := LastIndex(slices.All(s), matcher)
	if !ok {
		return IndexUnavailable
	}

	return i
}

// LastIndexOf returns the last index of the element and true.
// If no element was found, it returns zero value of K and false.
func LastIndexOf[K any, V comparable](seq iter.Seq2[K, V], expected V) (K, bool) {
	return LastIndex(seq, MatchEqual(expected))
}

// SliceLastIndexOf searches slice for last element which is equal to expected and returns index of it
// or IndexUnavailable if not found.
func SliceLastIndexOf[Slice ~[]E, E comparable](expected E, s Slice) int {
	return SliceLastIndex(s, MatchEqual(expected))
}

// Reduce reduces the seq to a single value using reducer to combine the elements. It performs reducer on all elements
// sequentially from first to last and returns output of last reducer call.
// Reducer accumulator is set to initial for first call.
// If the seq is empty, it returns initial value.
func Reduce[E, A any](seq iter.Seq[E], reducer Reducer[A, E], initial A) A {
	accum := initial
	for e := range seq {
		accum = reducer.Reduce(accum, e)
	}

	return accum
}

// SliceReduce takes slice s, performs reducer on all elements sequentially from start to end and returns output of last
// reducer call. Reducer accumulator is set to initial for first call. If the seq is empty, it returns initial value.
func SliceReduce[Slice ~[]E, E, A any](s Slice, reducer Reducer[A, E], initial A) A {
	return Reduce(slices.Values(s), reducer, initial)
}

// Filter returns a new iterator with elements that satisfy the predicate.
// The order of the elements is preserved.
// If no element satisfies the predicate, it returns an empty iterator.
func Filter[E any](seq iter.Seq[E], matcher Matcher[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range seq {
			if matcher.Match(e) {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// SliceFilter returns a new slice consisting of original slice elements which matched using matcher.
// It returns nil if no element was matched.
func SliceFilter[Slice ~[]E, E any](s Slice, matcher Matcher[E]) Slice {
	return slices.Collect(Filter(slices.Values(s), matcher))
}

// OrganicSliceFilter is like SliceFilter with more performance as it's implemented organic.
func OrganicSliceFilter[Slice ~[]E, E any](s Slice, matcher Matcher[E]) Slice {
	filtered := make(Slice, 0, len(s))
	for _, e := range s {
		if matcher.Match(e) {
			filtered = append(filtered, e)
		}
	}

	if len(filtered) == 0 {
		return nil
	}

	return filtered[:len(filtered):len(filtered)]
}

// Map returns a new seq with elements that are transformed by the transformer.
// Input and output sequences have same length.
func Map[S, D any](seq iter.Seq[S], transformer Transformer[S, D]) iter.Seq[D] {
	return func(yield func(D) bool) {
		for e := range seq {
			if !yield(transformer(e)) {
				return
			}
		}
	}
}

// SliceMap iterates over ss, applies transformer on each individual element of it and returns the transformed slice.
// Input and output slices have same length.
func SliceMap[SS ~[]S, DD ~[]D, S, D any](ss SS, transformer Transformer[S, D]) DD {
	return slices.Collect(Map(slices.Values(ss), transformer))
}

// OrganicSliceMap is like SliceMap with more performance as it's implemented organic.
func OrganicSliceMap[SS ~[]S, DD ~[]D, S, D any](ss SS, transformer Transformer[S, D]) DD {
	if ss == nil {
		return nil
	}

	dd := make(DD, len(ss))
	for k, src := range ss {
		dd[k] = transformer(src)
	}

	return dd
}

// FlatMap returns a new seq with elements that are transformed by the transformer and flattened into new sequence.
func FlatMap[S, D any](seq iter.Seq[S], transformer Transformer[S, iter.Seq[D]]) iter.Seq[D] {
	return func(yield func(D) bool) {
		for e := range seq {
			out := transformer(e)
			for d := range out {
				if !yield(d) {
					return
				}
			}
		}
	}
}

func SliceFlatMap[SS ~[]S, DD ~[]D, S, D any](ss SS, transformer Transformer[S, DD]) DD {
	return slices.Collect(FlatMap(slices.Values(ss), func(src S) iter.Seq[D] {
		return func(yield func(D) bool) {
			out := transformer(src)
			for _, d := range out {
				if !yield(d) {
					return
				}
			}
		}
	}))
}

// Flat merges multiple sequences of same type into one sequence. It does not modify nor drop elements. So output
// may contain elements with duplicate values from different sequences. It's an alternative to slices.Concat.
func Flat[E any](ss ...iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, s := range ss {
			for e := range s {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// SliceFlat merges multiple slices of same type(Slice) into one slice. It does not modify nor drop elements. So output
// may contain elements with duplicate values from different slices. It's an alternative to slices.Concat.
func SliceFlat[Slice ~[]E, E any](ss ...Slice) Slice {
	length := SliceReduce(ss, ReducerFunc[int, Slice](func(r int, s Slice) int {
		return r + len(s)
	}), 0)

	if length == 0 {
		return nil
	}

	flat := make(Slice, 0, length)
	for _, s := range ss {
		flat = append(flat, s...)
	}

	return flat
}

// DistinctBy returns a new seq with distinct elements.
// Transformer is used to determine the key of the element to check elements' equality.
func DistinctBy[E any, K comparable](seq iter.Seq[E], transformer Transformer[E, K]) iter.Seq[E] {
	return func(yield func(E) bool) {
		history := make(map[K]bool)

		for e := range seq {
			id := transformer(e)
			if _, ok := history[id]; !ok {
				history[id] = true

				if !yield(e) {
					return
				}
			}
		}
	}
}

// SliceDistinctBy performs DistinctBy on slice
func SliceDistinctBy[Slice ~[]E, E any, K comparable](s Slice, transformer Transformer[E, K]) Slice {
	return slices.Collect(DistinctBy(slices.Values(s), transformer))
}

// Distinct returns a new seq with distinct elements.
func Distinct[E comparable](seq iter.Seq[E]) iter.Seq[E] {
	return DistinctBy(seq, Identity)
}

// SliceDistinct performs Distinct on slice
func SliceDistinct[Slice ~[]E, E comparable](s Slice) Slice {
	return slices.Collect(Distinct(slices.Values(s)))
}

// First returns the first element in the iterator.
// If the iterator is empty, it returns the zero value of T and false.
func First[E any](seq iter.Seq[E]) (E, bool) {
	next, stop := iter.Pull(seq)
	defer stop()

	return next()
}

// Last returns the last element in the iterator.
// If the iterator is empty, it returns the zero value of T and false.
func Last[E any](seq iter.Seq[E]) (E, bool) {
	var v E
	ok := false

	for e := range seq {
		v, ok = e, true
	}

	return v, ok
}

// ElementAt returns the element at the specified index.
// If the index is out of range, it returns the zero value of E and false.
func ElementAt[E any](seq iter.Seq[E], index int) (E, bool) {
	var v E
	if index < 0 {
		return v, false
	}

	i := 0
	for e := range seq {
		if i == index {
			return e, true
		}

		i++
	}

	return v, false
}

// ElementAtOrDefault returns the element at the specified index.
// If the index is out of range, it returns the default value.
func ElementAtOrDefault[E any](seq iter.Seq[E], index int, defaultValue E) E {
	v, ok := ElementAt(seq, index)
	if !ok {
		return defaultValue
	}

	return v
}

// ElementAtOrElse returns the element at the specified index.
// If the index is out of range, it returns the value from the defaultFunc.
func ElementAtOrElse[E any](seq iter.Seq[E], index int, defaultFunc func() E) E {
	v, ok := ElementAt(seq, index)
	if !ok {
		return defaultFunc()
	}

	return v
}

// Count returns the number of elements in the seq.
func Count[E any](seq iter.Seq[E]) int {
	c := 0
	for range seq {
		c++
	}

	return c
}

// ForEach calls the function for each element in the seq.
func ForEach[E any](seq iter.Seq[E], f func(E)) {
	for e := range seq {
		f(e)
	}
}

// Skip returns a new seq with the first n elements skipped.
func Skip[E any](seq iter.Seq[E], n int) iter.Seq[E] {
	return func(yield func(E) bool) {
		i := 0
		for e := range seq {
			if i >= n {
				if !yield(e) {
					return
				}
			}

			i++
		}
	}
}

// SkipWhile returns a new seq with elements skipped while the predicate is true.
func SkipWhile[E any](seq iter.Seq[E], matcher Matcher[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		skipping := true
		for e := range seq {
			if skipping && matcher.Match(e) {
				continue
			}

			skipping = false
			if !yield(e) {
				return
			}
		}
	}
}

// Take returns a new seq with the first n elements taken.
func Take[E any](seq iter.Seq[E], n int) iter.Seq[E] {
	return func(yield func(E) bool) {
		i := 0
		for e := range seq {
			if i >= n || !yield(e) {
				return
			}

			i++
		}
	}
}

// TakeWhile returns a new seq with elements taken while the predicate is true.
func TakeWhile[E any](seq iter.Seq[E], matcher Matcher[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for e := range seq {
			if !matcher.Match(e) || !yield(e) {
				return
			}
		}
	}
}

// Chain returns a new seq that chains multiple sequences.
func Chain[E any](ss ...iter.Seq[E]) iter.Seq[E] {
	return func(yield func(E) bool) {
		for _, seq := range ss {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Extend returns a new seq that chains a seq with a slice.
func Extend[Slice ~[]E, E any](seq iter.Seq[E], s Slice) iter.Seq[E] {
	return Chain(seq, slices.Values(s))
}
