package conv

func PtrVal[T any](src *T) T {
	if src == nil {
		var dst T

		return dst
	}

	return *src
}

func ValPtr[T any](src T) *T {
	return &src
}
