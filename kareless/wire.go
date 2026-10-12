package kareless

// Route is a destination on one medium. Medium is caller-defined. Address is the medium-specific address.
type Route struct {
	// Medium identifies the transport. Callers choose the numbering.
	Medium int
	// Address is the medium-specific destination.
	Address string
}

// Router maps a logical address to a Route.
type Router interface {
	// Resolve returns the route for addr.
	Resolve(addr string) Route
}

// RouterFunc adapts a function to Router.
type RouterFunc func(addr string) Route

// Resolve calls f.
func (f RouterFunc) Resolve(addr string) Route {
	return f(addr)
}

// Marshaler encodes a payload into bytes.
type Marshaler interface {
	// Marshal encodes payload.
	Marshal(payload any) []byte
}

// MarshalerFunc adapts a function to Marshaler.
type MarshalerFunc func(payload any) []byte

// Marshal calls f.
func (f MarshalerFunc) Marshal(payload any) []byte {
	return f(payload)
}

// Unmarshaler decodes bytes into v.
type Unmarshaler interface {
	// Unmarshal decodes data into v.
	Unmarshal(data []byte, v any) error
}

// UnmarshalerFunc adapts a function to Unmarshaler.
type UnmarshalerFunc func(data []byte, v any) error

// Unmarshal calls f.
func (f UnmarshalerFunc) Unmarshal(data []byte, v any) error {
	return f(data, v)
}

// Encapsulator binds encoded bytes to a route and returns the message that goes on the wire.
type Encapsulator[M any] interface {
	// Encapsulate binds data to route.
	Encapsulate(route Route, data []byte) M
}

// EncapsulatorFunc adapts a function to Encapsulator.
type EncapsulatorFunc[M any] func(route Route, data []byte) M

// Encapsulate calls f.
func (f EncapsulatorFunc[M]) Encapsulate(route Route, data []byte) M {
	return f(route, data)
}

// Decapsulator extracts encoded bytes from a message taken off the wire.
type Decapsulator[M any] interface {
	// Decapsulate extracts the payload bytes from msg.
	Decapsulate(msg M) ([]byte, error)
}

// DecapsulatorFunc adapts a function to Decapsulator.
type DecapsulatorFunc[M any] func(msg M) ([]byte, error)

// Decapsulate calls f.
func (f DecapsulatorFunc[M]) Decapsulate(msg M) ([]byte, error) {
	return f(msg)
}

// Muldem is a bidirectional multiplexer and demultiplexer which is able to:
//  1. Marshal message using the Marshaler and encapsulate it for a specific address
//     using the Encapsulator and Router to be put on the wire.
//  2. Decapsulate message picked from the wire from a specific address
//     using the Decapsulator and unmarshal it using the Unmarshaler.
type Muldem[M any] struct {
	// Router maps a logical address to a Route.
	Router Router

	// Marshaler encodes payloads. Unmarshaler decodes them.
	Marshaler   Marshaler
	Unmarshaler Unmarshaler

	// Encapsulator binds encoded bytes to a route. Decapsulator extracts them again.
	Encapsulator Encapsulator[M]
	Decapsulator Decapsulator[M]
}

// WithRouter returns a copy that uses r.
func (mx Muldem[M]) WithRouter(r Router) Muldem[M] {
	mx.Router = r

	return mx
}

// WithEncoding returns a copy that uses m and um.
func (mx Muldem[M]) WithEncoding(m Marshaler, um Unmarshaler) Muldem[M] {
	mx.Marshaler = m
	mx.Unmarshaler = um

	return mx
}

// WithEncapsulation returns a copy that uses e and de.
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
