package tricks

func Copy[T any](tt []T) []T {
	ftt := make([]T, len(tt))
	copy(ftt, tt)

	return ftt
}
