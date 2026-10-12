package std

import (
	"encoding/json"
	"errors"

	"github.com/janstoon/toolbox/bricks"
	"google.golang.org/protobuf/proto"

	"github.com/janstoon/toolbox/kareless"
)

// IdentityRouter returns a router that keeps addr unchanged and tags every route with medium.
func IdentityRouter(medium int) kareless.Router {
	return kareless.RouterFunc(func(addr string) kareless.Route {
		return kareless.Route{
			Medium:  medium,
			Address: addr,
		}
	})
}

// JsonMarshaler encodes payloads with encoding/json. A marshal error yields nil.
var JsonMarshaler = kareless.MarshalerFunc(func(payload any) []byte {
	bb, _ := json.Marshal(payload)

	return bb
})

// JsonUnmarshaler decodes JSON with encoding/json.
var JsonUnmarshaler = kareless.UnmarshalerFunc(json.Unmarshal)

// ProtobufMarshaler encodes a proto.Message. A value that is not a proto.Message yields nil.
var ProtobufMarshaler = kareless.MarshalerFunc(func(payload any) []byte {
	msg, ok := payload.(proto.Message)
	if !ok {
		return nil
	}

	bb, _ := proto.Marshal(msg)

	return bb
})

// ProtobufUnmarshaler decodes bytes into a proto.Message.
// It returns an error that wraps bricks.ErrInvalidArgument when v is not a proto.Message.
var ProtobufUnmarshaler = kareless.UnmarshalerFunc(func(bb []byte, v any) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return errors.Join(bricks.ErrInvalidArgument, errors.New("variable is not proto.Message"))
	}

	return proto.Unmarshal(bb, msg)
})

// NoopEncapsulator returns an encapsulator that discards the route and the bytes
// and returns the zero message.
func NoopEncapsulator[M any]() kareless.Encapsulator[M] {
	return kareless.EncapsulatorFunc[M](func(route kareless.Route, data []byte) (m M) {
		return m
	})
}

// NoopDecapsulator returns a decapsulator that discards the message and returns empty bytes.
func NoopDecapsulator[M any]() kareless.Decapsulator[M] {
	return kareless.DecapsulatorFunc[M](func(msg M) ([]byte, error) {
		return nil, nil
	})
}

// NewMuldem creates a simple kareless.Muldem[M] with IdentityRouter(0), Json Encoding (JsonMarshaler, JsonUnmarshaler)
// and Noop Encapsulation (NoopEncapsulator[M], NoopDecapsulator[M]).
func NewMuldem[M any]() kareless.Muldem[M] {
	return kareless.Muldem[M]{
		Router: IdentityRouter(0),

		Marshaler:   JsonMarshaler,
		Unmarshaler: JsonUnmarshaler,

		Encapsulator: NoopEncapsulator[M](),
		Decapsulator: NoopDecapsulator[M](),
	}
}
