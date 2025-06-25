package tricks

// Transformer is a method which converts S to D
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

// ToAny casts type T to any. One use-case is to feed a method which accepts []any, but you have []T in hand
func ToAny[T any](src T) any {
	return src
}

// FromAny tries to cast any to type T. If cast failed it returns zero value of type T.
// One use-case is to feed a method which accepts []T, but you have []ant in hand
func FromAny[T any](src any) T {
	dst, _ := src.(T)

	return dst
}
