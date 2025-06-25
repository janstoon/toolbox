package tricks

// Option is a modifier of S. It applies a modification on S and return it.
type Option[S any] interface {
	Apply(s *S) *S
}

type (
	// InPlaceOption takes S reference and modifies it directly
	InPlaceOption[S any] func(s *S)

	// OutOfPlaceOption takes a copy of S and returns the modified copy of it
	OutOfPlaceOption[S any] func(s S) S
)

func (o InPlaceOption[S]) Apply(s *S) *S {
	o(s)

	return s
}

func (o OutOfPlaceOption[S]) Apply(s *S) *S {
	return ValPtr(o(*s))
}

// ApplyOptions performs modifications of all Options on S
func ApplyOptions[S any](s *S, oo ...Option[S]) *S {
	d := s
	for _, opt := range oo {
		d = opt.Apply(d)
	}

	return d
}
