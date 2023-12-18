package tricks

// Coalesce returns left-most non-zero value
func Coalesce[T comparable](tt ...T) T {
	var zero T
	if len(tt) == 0 {
		return zero
	}

	if tt[0] != zero || len(tt) == 1 {
		return tt[0]
	}

	return Coalesce(tt[1:]...)
}
