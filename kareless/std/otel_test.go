package std_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janstoon/toolbox/kareless"
	"github.com/janstoon/toolbox/kareless/std"
)

func TestOpenTelemetryDriverConstructor(t *testing.T) {
	ss := new(kareless.Settings)
	ss.Append(std.MapSettingSource{
		std.TelemetryCollectorOtlpHostSettingKey:            "127.0.0.1",
		std.TelemetryCollectorOtlpPortsSettingKey + ".grpc": 4317,
		std.TelemetryCollectorOtlpPortsSettingKey + ".http": 4318,
	})

	drv := std.OpenTelemetryDriverConstructor("toolbox")(ss, nil, nil)
	assert.NotNil(t, drv)
}
