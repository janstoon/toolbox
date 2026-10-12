package std_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/kareless"
	"github.com/janstoon/toolbox/kareless/std"
)

func TestIdentityRouterAndCodecs(t *testing.T) {
	route := std.IdentityRouter(7).Resolve("a")
	assert.Equal(t, kareless.Route{Medium: 7, Address: "a"}, route)

	raw := std.JsonMarshaler.Marshal(map[string]string{"k": "v"})
	var decoded map[string]string
	require.NoError(t, std.JsonUnmarshaler.Unmarshal(raw, &decoded))
	assert.Equal(t, "v", decoded["k"])
	assert.Nil(t, std.JsonMarshaler.Marshal(func() {}))

	assert.Nil(t, std.ProtobufMarshaler.Marshal("nope"))
	require.ErrorIs(t, std.ProtobufUnmarshaler.Unmarshal(nil, "nope"), bricks.ErrInvalidArgument)

	msg := wrapperspb.String("hi")
	raw = std.ProtobufMarshaler.Marshal(msg)
	out := wrapperspb.String("")
	require.NoError(t, std.ProtobufUnmarshaler.Unmarshal(raw, out))
	assert.Equal(t, "hi", out.GetValue())

	assert.Zero(t, std.NoopEncapsulator[int]().Encapsulate(route, []byte("x")))
	body, err := std.NoopDecapsulator[int]().Decapsulate(1)
	require.NoError(t, err)
	assert.Nil(t, body)

	assert.Zero(t, std.NewMuldem[int]().Encapsulate("a", map[string]int{"n": 1}))
}

func ExampleNewMuldem() {
	mx := std.NewMuldem[int]()
	fmt.Println(mx.Encapsulate("ignored", map[string]string{"k": "v"}))

	// Output:
	// 0
}
