package tricks

// Transformer is any method which converts S to D. It's usable in slice elements transformations, casting and
// interlayer data structure mapping like mapping entity to dao or dto and vice versa.
type Transformer[S, D any] func(src S) D

// Identity is a transformer which returns input without any modification.
// It's also the simplest Middleware which directly returns the next handler aka passthrough middleware.
func Identity[T any](t T) T { return t }

// PtrVal dereferences a pointer of type T. If it's nil it returns the zero value of T.
// It's a safe pointer dereference method intended for use cases where a value is required,
// even if the pointer is nil or empty.
func PtrVal[T any](src *T) T {
	if src == nil {
		var dst T

		return dst
	}

	return *src
}

// ValPtr converts a value of type T to its pointer. It simply returns address of src and is optimal if inlined by
// the compiler. Since input is passed by value it returns a pointer to the copied variable and further modifications
// on it does not change value of the input variable. This method is intended for use cases where expected type is a
// pointer for reasons other than modification like protobuf data structures, api dto or storage dao.
func ValPtr[T any](src T) *T {
	return &src
}

// PtrPtr transforms pointer of type S to pointer of type D using Transformer. If src is nil output is nil too.
// It's a safe transformation with nil guard and safe pointer dereference.
func PtrPtr[S, D any](src *S, transformer Transformer[S, D]) *D {
	if src == nil {
		return nil
	}

	return ValPtr(transformer(PtrVal(src)))
}

// ToAny casts type T to any. One use-case is to feed a method which accepts []any, but you have []T in hand.
func ToAny[T any](src T) any {
	return src
}

// FromAny tries to cast any to type T. If cast failed it returns zero value of type T.
// One use-case is to feed a method which accepts []T, but you have []ant in hand.
func FromAny[T any](src any) T {
	dst, _ := src.(T)

	return dst
}

// StringToRunes casts string into slice of underlying rune(s)
func StringToRunes(src string) []rune {
	return []rune(src)
}
