package kareless_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/janstoon/toolbox/kareless"
	"github.com/janstoon/toolbox/kareless/std"
)

type driverFunc func(context.Context) error

func (f driverFunc) Run(ctx context.Context) error { return f(ctx) }

func TestCompileOptions(t *testing.T) {
	hookRan := make(chan struct{})
	k := kareless.Compile(
		kareless.Feeder(std.MapSettingSource{"name": "ada"}),
		kareless.Equipment(func(ss *kareless.Settings) []kareless.InstrumentCatalogue {
			return []kareless.InstrumentCatalogue{{
				Names: []string{"name"},
				Builder: func(ss *kareless.Settings, _ *kareless.InstrumentBank) kareless.Instrument {
					return ss.GetString("name")
				},
			}}
		}),
		kareless.Installer(func(ss *kareless.Settings, ib *kareless.InstrumentBank) kareless.Application {
			return kareless.ResolveInstrumentByType[string](ib, "name")
		}),
		kareless.Connector(func(_ *kareless.Settings, _ *kareless.InstrumentBank, apps []kareless.Application) kareless.Driver {
			return driverFunc(func(context.Context) error {
				if len(apps) != 1 || apps[0].(string) != "ada" {
					return assert.AnError
				}

				return nil
			})
		}),
		kareless.PostHook(func(_ context.Context, ss *kareless.Settings, ib *kareless.InstrumentBank, _ []kareless.Application) error {
			assert.Equal(t, "ada", ss.GetString("name"))
			assert.Equal(t, "ada", kareless.ResolveInstrumentByType[string](ib, "name"))
			close(hookRan)

			return nil
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- k.Run(ctx)
	}()

	select {
	case <-hookRan:
	case <-time.After(2 * time.Second):
		t.Fatal("after-start hook did not run")
	}

	require.NoError(t, <-errCh)
}

func TestSettingsGetters(t *testing.T) {
	ss := new(kareless.Settings)
	ss.Prepend(std.MapSettingSource{"name": "fallback"})
	ss.Append(std.MapSettingSource{
		"name":   "ada",
		"tags":   []any{"a", "b"},
		"n":      3,
		"ok":     true,
		"wait":   "2s",
		"nested": map[string]any{"city": "Shiraz"},
		"list":   []any{"x"},
	})

	assert.Equal(t, "fallback", ss.GetString("name"))
	assert.Equal(t, []string{"a", "b"}, ss.GetStringSlice("tags"))
	assert.Equal(t, 3, ss.GetInt("n"))
	assert.Equal(t, int64(3), ss.GetInt64("n"))
	assert.Equal(t, byte(3), ss.GetByte("n"))
	assert.True(t, ss.GetBool("ok"))
	assert.Equal(t, 2*time.Second, ss.GetDuration("wait"))
	assert.ElementsMatch(t, []string{"city"}, ss.Children("nested"))
	assert.Equal(t, []string{"0"}, ss.Children("list"))
	assert.Nil(t, ss.Children("missing"))
	assert.Empty(t, ss.GetString("missing"))
	assert.False(t, ss.GetBool("missing"))
	assert.Zero(t, ss.GetDuration("missing"))
}
