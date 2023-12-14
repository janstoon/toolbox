package tricks

type Transformer[S, D any] func(src S) D

// PtrVal dereferences a pointer of type T. If it's nil it returns the zero value of T
func PtrVal[T any](src *T) T {
	if src == nil {
		var dst T

		return dst
	}

	return *src
}

// ValPtr converts a value of type T to its pointer
func ValPtr[T any](src T) *T {
	return &src
}

// PtrPtr transforms pointer of type S to pointer of type D. If src is nil it returns nil
func PtrPtr[S, D any](transformer Transformer[S, D], src *S) *D {
	if src == nil {
		return nil
	}

	return ValPtr(transformer(PtrVal(src)))
}
