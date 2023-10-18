package tricks

type Option[S any] interface {
	Apply(s *S) *S
}

type (
	InPlaceOption[S any]    func(s *S)
	OutOfPlaceOption[S any] func(s S) S
)

func (o InPlaceOption[S]) Apply(s *S) *S {
	o(s)

	return s
}

func (o OutOfPlaceOption[S]) Apply(s *S) *S {
	return ValPtr(o(*s))
}

func ApplyOptions[S any](s *S, oo ...Option[S]) *S {
	d := s
	for _, opt := range oo {
		d = opt.Apply(d)
	}

	return d
}
