package tricks

func Filter[T any](match func(src T) bool, tt []T) []T {
	ftt := make([]T, 0, len(tt))
	for _, t := range tt {
		if match(t) {
			ftt = append(ftt, t)
		}
	}

	return ftt[:len(ftt):len(ftt)]
}

func Map[S, D any](transformer Transformer[S, D], ss []S) []D {
	dd := make([]D, len(ss))
	for k, src := range ss {
		dd[k] = transformer(src)
	}

	return dd
}

func Reduce[T, R any](reducer func(r R, t T) R, tt []T) R {
	var r R
	for _, t := range tt {
		r = reducer(r, t)
	}

	return r
}

func Find[T any](match func(src T) bool, tt []T) *T {
	for _, t := range tt {
		if match(t) {
			return ValPtr(t)
		}
	}

	return nil
}

func FindIndex[T any](match func(src T) bool, tt []T) int {
	for i, t := range tt {
		if match(t) {
			return i
		}
	}

	return -1
}

func IndexOf[T comparable](expected T, tt []T) int {
	return FindIndex(MatchEqual(expected), tt)
}

func Flat[T any](slices ...[]T) []T {
	length := Reduce(func(r int, tt []T) int {
		return r + len(tt)
	}, slices)

	flat := make([]T, 0, length)
	for _, slice := range slices {
		flat = append(flat, slice...)
	}

	return flat
}
