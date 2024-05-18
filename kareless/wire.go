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
	Unmarshal(data []byte) any
}

type UnmarshalerFunc func(data []byte) any

func (f UnmarshalerFunc) Unmarshal(data []byte) any {
	return f(data)
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
