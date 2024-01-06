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
func PtrPtr[S, D any](src *S, transformer Transformer[S, D]) *D {
	if src == nil {
		return nil
	}

	return ValPtr(transformer(PtrVal(src)))
}

func ToAny[T any](src T) any {
	return src
}

func FromAny[T any](src any) T {
	dst, _ := src.(T)

	return dst
}
