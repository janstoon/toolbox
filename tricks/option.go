package tricks

// Option is a modifier of S. It applies a modification on S and return it.
// It is intended to perform partial and isolated updates using ApplyOptions.
// One usage is functional options pattern that is instead of passing a config struct or many arguments
// to a constructor, you pass a variadic list of functions (options) that modify the internal state.
// The modifier can implement this interface or can be a function and used as an Option
// with help of MutableOption and ImmutableOption.
//
// Another use-case is in making object Transformer(s) with optional field modifiers.
//
// The pattern was popularized by Dave Cheney and Rob Pike.
// For more details visit https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
//
// Example:
//
//	type Server struct {
//		Host string
//		Port int
//		TLS  bool
//	}
//
//	// Constructor with functional options
//	func NewServer(opts ...Option[Server]) *Server {
//		s := &Server{
//			Host: "localhost", // default values
//			Port: 8080,
//			TLS:  false,
//		}
//
//		return ApplyOptions(s, opts...)
//	}
//
//	// Individual options
//
//	func WithHost(host string) Option[Server] {
//		return MutableOption[Server](func(s *Server) {
//			s.Host = host
//		})
//	}
//
//	func WithPort(port int) Option[Server] {
//		return ImmutableOption[Server](func(s Server) Server {
//			s.Port = port
//
//			return s
//		})
//	}
//
//	func WithTLS(enabled bool) Option[Server] {
//		return ImmutableOption[Server](func(s Server) Server {
//			s.TLS = enabled
//
//			return s
//		})
//	}
//
// Usage:
//
//	s := NewServer(
//		WithHost("example.com"),
//		WithPort(443),
//		WithTLS(true),
//	)
type Option[S any] interface {
	Apply(s *S) *S
}

// MutableOption takes S reference and modifies it directly.
type MutableOption[S any] func(s *S)

func (o MutableOption[S]) Apply(s *S) *S {
	o(s)

	return s
}

// ImmutableOption takes a copy of S and returns the modified copy of it.
type ImmutableOption[S any] func(s S) S

func (o ImmutableOption[S]) Apply(s *S) *S {
	return ValPtr(o(*s))
}

// ApplyOptions performs modifications of all Option(s) on S sequentially.
func ApplyOptions[S any](s *S, oo ...Option[S]) *S {
	d := s
	for _, opt := range oo {
		d = opt.Apply(d)
	}

	return d
}
