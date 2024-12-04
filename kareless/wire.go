package kareless

type Route struct {
	Medium  int
	Address string // Medium-specific address
}

type Router interface {
	Resolve(addr string) Route
}

type RouterFunc func(addr string) Route

func (f RouterFunc) Resolve(addr string) Route {
	return f(addr)
}

type Marshaler interface {
	Marshal(payload any) []byte
}

type MarshalerFunc func(payload any) []byte

func (f MarshalerFunc) Marshal(payload any) []byte {
	return f(payload)
}

type Unmarshaler interface {
	Unmarshal(data []byte, v any) error
}

type UnmarshalerFunc func(data []byte, v any) error

func (f UnmarshalerFunc) Unmarshal(data []byte, v any) error {
	return f(data, v)
}

type Encapsulator[M any] interface {
	Encapsulate(route Route, data []byte) M
}

type EncapsulatorFunc[M any] func(route Route, data []byte) M

func (f EncapsulatorFunc[M]) Encapsulate(route Route, data []byte) M {
	return f(route, data)
}

type Decapsulator[M any] interface {
	Decapsulate(msg M) ([]byte, error)
}

type DecapsulatorFunc[M any] func(msg M) ([]byte, error)

func (f DecapsulatorFunc[M]) Decapsulate(msg M) ([]byte, error) {
	return f(msg)
}

// Muldem is a bidirectional multiplexer and demultiplexer which is able to:
//  1. Marshal message using the Marshaler and encapsulate it for a specific address
//     using the Encapsulator and Router to be put on the wire.
//  2. Decapsulate message picked from the wire from a specific address
//     using the Decapsulator and unmarshal it using the Unmarshaler.
type Muldem[M any] struct {
	Router Router

	Marshaler   Marshaler
	Unmarshaler Unmarshaler

	Encapsulator Encapsulator[M]
	Decapsulator Decapsulator[M]
}

func (mx Muldem[M]) WithRouter(r Router) Muldem[M] {
	mx.Router = r

	return mx
}

func (mx Muldem[M]) WithEncoding(m Marshaler, um Unmarshaler) Muldem[M] {
	mx.Marshaler = m
	mx.Unmarshaler = um

	return mx
}

func (mx Muldem[M]) WithEncapsulation(e Encapsulator[M], de Decapsulator[M]) Muldem[M] {
	mx.Encapsulator = e
	mx.Decapsulator = de

	return mx
}

// Encapsulate marshals the payload and binds it to a specific address using Router
// and outputs the Message(M) ready to be put on the wire.
func (mx Muldem[M]) Encapsulate(addr string, payload any) M {
	return mx.Encapsulator.Encapsulate(mx.Router.Resolve(addr), mx.Marshaler.Marshal(payload))
}
