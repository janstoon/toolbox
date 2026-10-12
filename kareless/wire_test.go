package kareless_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/janstoon/toolbox/kareless"
)

func TestMuldem(t *testing.T) {
	type msg struct {
		Addr string
		Body []byte
	}

	mx := kareless.Muldem[msg]{}.
		WithRouter(kareless.RouterFunc(func(addr string) kareless.Route {
			return kareless.Route{Medium: 1, Address: addr}
		})).
		WithEncoding(
			kareless.MarshalerFunc(func(payload any) []byte {
				return []byte(payload.(string))
			}),
			kareless.UnmarshalerFunc(func(data []byte, v any) error {
				*v.(*string) = string(data)

				return nil
			}),
		).
		WithEncapsulation(
			kareless.EncapsulatorFunc[msg](func(route kareless.Route, data []byte) msg {
				return msg{Addr: route.Address, Body: data}
			}),
			kareless.DecapsulatorFunc[msg](func(m msg) ([]byte, error) {
				return m.Body, nil
			}),
		)

	got := mx.Encapsulate("inbox", "hi")
	assert.Equal(t, msg{Addr: "inbox", Body: []byte("hi")}, got)

	var decoded string
	require.NoError(t, mx.Unmarshaler.Unmarshal(got.Body, &decoded))
	assert.Equal(t, "hi", decoded)
}

func ExampleMuldem_Encapsulate() {
	type msg struct {
		Addr string
		Body []byte
	}

	mx := kareless.Muldem[msg]{}.
		WithRouter(kareless.RouterFunc(func(addr string) kareless.Route {
			return kareless.Route{Address: addr}
		})).
		WithEncoding(
			kareless.MarshalerFunc(func(payload any) []byte { return []byte(fmt.Sprint(payload)) }),
			kareless.UnmarshalerFunc(func([]byte, any) error { return nil }),
		).
		WithEncapsulation(
			kareless.EncapsulatorFunc[msg](func(route kareless.Route, data []byte) msg {
				return msg{Addr: route.Address, Body: data}
			}),
			kareless.DecapsulatorFunc[msg](func(msg) ([]byte, error) { return nil, nil }),
		)

	fmt.Printf("%s %s\n", mx.Encapsulate("inbox", "hi").Addr, mx.Encapsulate("inbox", "hi").Body)

	// Output:
	// inbox hi
}
