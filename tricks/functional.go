package tricks

func Map[S, D any](ss []S, transformer Transformer[S, D]) []D {
	if ss == nil {
		return nil
	}

	dd := make([]D, len(ss))
	for k, src := range ss {
		dd[k] = transformer(src)
	}

	return dd
}

func Reduce[T, R any](tt []T, reducer func(r R, t T) R) R {
	var r R
	for _, t := range tt {
		r = reducer(r, t)
	}

	return r
}

func Filter[T any](tt []T, match func(src T) bool) []T {
	ftt := make([]T, 0, len(tt))
	for _, t := range tt {
		if match(t) {
			ftt = append(ftt, t)
		}
	}

	if len(ftt) == 0 {
		return nil
	}

	return ftt[:len(ftt):len(ftt)]
}

func Find[T any](tt []T, match func(src T) bool) *T {
	for _, t := range tt {
		if match(t) {
			return ValPtr(t)
		}
	}

	return nil
}

const IndexUnavailable = -1

func FindIndex[T any](tt []T, match func(src T) bool) int {
	for i, t := range tt {
		if match(t) {
			return i
		}
	}

	return IndexUnavailable
}

func IndexOf[T comparable](expected T, tt []T) int {
	return FindIndex(tt, MatchEqual(expected))
}

func Flat[T any](slices ...[]T) []T {
	length := Reduce(slices, func(r int, tt []T) int {
		return r + len(tt)
	})

	if length == 0 {
		return nil
	}

	flat := make([]T, 0, length)
	for _, slice := range slices {
		flat = append(flat, slice...)
	}

	return flat
}
