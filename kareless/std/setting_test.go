package std_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/janstoon/toolbox/bricks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/janstoon/toolbox/kareless/std"
)

func TestMapSettingSource_Get(t *testing.T) {
	ctx := context.Background()

	ss := std.MapSettingSource{
		"c1": "v-c1",
		"n1": map[string]map[string]string{
			"a": {
				"b": "v-n1ab",
			},
		},
		"n1.w": "v-n1w",
	}

	v, err := ss.Get(ctx, "")
	require.NoError(t, err)
	assert.Equal(t, std.MapSettingSource{
		"c1": "v-c1",
		"n1": map[string]map[string]string{
			"a": {
				"b": "v-n1ab",
			},
		},
		"n1.w": "v-n1w",
	}, v)

	v, err = ss.Get(ctx, "c1")
	require.NoError(t, err)
	assert.Equal(t, "v-c1", v)

	v, err = ss.Get(ctx, "n1")
	require.NoError(t, err)
	assert.Equal(t, map[string]map[string]string{
		"a": {
			"b": "v-n1ab",
		},
	}, v)

	v, err = ss.Get(ctx, "n1.a")
	require.NoError(t, err)
	assert.Equal(t, map[string]any{
		"b": "v-n1ab",
	}, v)

	v, err = ss.Get(ctx, "n1.a.b")
	require.NoError(t, err)
	assert.Equal(t, "v-n1ab", v)

	v, err = ss.Get(ctx, "n1.w")
	require.NoError(t, err)
	assert.Equal(t, "v-n1w", v)

	v, err = ss.Get(ctx, "unknown")
	require.ErrorIs(t, err, bricks.ErrNotFound)
	assert.Empty(t, v)

	v, err = ss.Get(ctx, "n1.unknown")
	require.ErrorIs(t, err, bricks.ErrNotFound)
	assert.Empty(t, v)

	v, err = ss.Get(ctx, "n1.a.unknown")
	require.ErrorIs(t, err, bricks.ErrNotFound)
	assert.Empty(t, v)
}

func TestLocalEarlySettings_Get(t *testing.T) {
	dir := t.TempDir()

	fname := "tconf"
	fh, err := os.Create(path.Join(dir, fmt.Sprintf("%s.json", fname)))
	require.NoError(t, err)
	defer func() { _ = fh.Close() }()

	enc := json.NewEncoder(fh)
	ctx := context.Background()

	ss := std.LocalEarlyLoadedSettingSource("")
	v, err := ss.Get(ctx, "c1")
	assert.Empty(t, v)
	require.ErrorIs(t, err, bricks.ErrNotFound)

	cfgOne := struct {
		C1 string `json:"C1"`
		C2 string `json:"c2"`
	}{
		C1: "Conf1ValueFromFile",
		C2: "Conf2ValueFromFile",
	}
	err = enc.Encode(cfgOne)
	require.NoError(t, err)
	err = fh.Sync()
	require.NoError(t, err)
	ss = std.LocalEarlyLoadedSettingSource(fname, dir)
	v, err = ss.Get(ctx, "")
	assert.Equal(t, map[string]any{"c1": "Conf1ValueFromFile", "c2": "Conf2ValueFromFile"}, v)
	require.NoError(t, err)
	v, err = ss.Get(ctx, ".")
	assert.Equal(t, map[string]any{"c1": "Conf1ValueFromFile", "c2": "Conf2ValueFromFile"}, v)
	require.NoError(t, err)
	v, err = ss.Get(ctx, "c1")
	assert.Equal(t, "Conf1ValueFromFile", v)
	require.NoError(t, err)
	v, err = ss.Get(ctx, "C1")
	assert.Equal(t, "Conf1ValueFromFile", v)
	require.NoError(t, err)
	v, err = ss.Get(ctx, "c2")
	assert.Equal(t, "Conf2ValueFromFile", v)
	require.NoError(t, err)
	v, err = ss.Get(ctx, "C2")
	assert.Equal(t, "Conf2ValueFromFile", v)
	require.NoError(t, err)

	t.Setenv("C1", "Conf1ValueFromEnv")
	ss = std.LocalEarlyLoadedSettingSource(fname, dir)
	v, err = ss.Get(ctx, "c1")
	assert.Equal(t, "Conf1ValueFromEnv", v)
	require.NoError(t, err)
	v, err = ss.Get(ctx, "C1")
	assert.Equal(t, "Conf1ValueFromEnv", v)
	require.NoError(t, err)
	v, err = ss.Get(ctx, "c2")
	assert.Equal(t, "Conf2ValueFromFile", v)
	require.NoError(t, err)
	v, err = ss.Get(ctx, "C2")
	assert.Equal(t, "Conf2ValueFromFile", v)
	require.NoError(t, err)

	err = os.Unsetenv("C1")
	require.NoError(t, err)
	ss = std.LocalEarlyLoadedSettingSource(fname, dir)
	v, err = ss.Get(ctx, "c1")
	assert.Equal(t, "Conf1ValueFromFile", v)
	require.NoError(t, err)
	v, err = ss.Get(ctx, "C1")
	assert.Equal(t, "Conf1ValueFromFile", v)
	require.NoError(t, err)
	v, err = ss.Get(ctx, "c2")
	assert.Equal(t, "Conf2ValueFromFile", v)
	require.NoError(t, err)
	v, err = ss.Get(ctx, "C2")
	assert.Equal(t, "Conf2ValueFromFile", v)
	require.NoError(t, err)

	cfgTwo := struct {
		C3 string `json:"c3"`
	}{
		C3: "Conf3ValueFromFile",
	}
	err = fh.Truncate(0)
	require.NoError(t, err)
	_, err = fh.Seek(0, 0)
	require.NoError(t, err)
	err = enc.Encode(cfgTwo)
	require.NoError(t, err)
	err = fh.Sync()
	require.NoError(t, err)
	ss = std.LocalEarlyLoadedSettingSource(fname, dir)
	v, err = ss.Get(ctx, "c1")
	assert.Empty(t, v)
	require.ErrorIs(t, err, bricks.ErrNotFound)
	v, err = ss.Get(ctx, "C1")
	assert.Empty(t, v)
	require.ErrorIs(t, err, bricks.ErrNotFound)
	v, err = ss.Get(ctx, "c2")
	assert.Empty(t, v)
	require.ErrorIs(t, err, bricks.ErrNotFound)
	v, err = ss.Get(ctx, "C2")
	assert.Empty(t, v)
	require.ErrorIs(t, err, bricks.ErrNotFound)
	v, err = ss.Get(ctx, "c3")
	assert.Equal(t, "Conf3ValueFromFile", v)
	require.NoError(t, err)
	v, err = ss.Get(ctx, "C3")
	assert.Equal(t, "Conf3ValueFromFile", v)
	require.NoError(t, err)
}
