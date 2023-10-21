package std_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/janstoon/toolbox/bricks"
	"github.com/janstoon/toolbox/tricks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/janstoon/toolbox/kareless/std"
)

func TestLocalSettings_Get(t *testing.T) {
	dir, err := os.MkdirTemp("", "kareless-std-local-settings-")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(dir) }()

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
		C1 string `json:"c1"`
	}{
		C1: "ValueFromFile",
	}
	err = enc.Encode(cfgOne)
	require.NoError(t, err)
	err = fh.Sync()
	require.NoError(t, err)
	ss = std.LocalEarlyLoadedSettingSource(fname, dir)
	v, err = ss.Get(ctx, "c1")
	assert.Equal(t, "ValueFromFile", tricks.PtrVal(v))
	require.NoError(t, err)

	t.Setenv("C1", "ValueFromEnv")
	ss = std.LocalEarlyLoadedSettingSource(fname, dir)
	v, err = ss.Get(ctx, "c1")
	assert.Equal(t, "ValueFromEnv", tricks.PtrVal(v))
	require.NoError(t, err)

	err = os.Unsetenv("C1")
	require.NoError(t, err)
	ss = std.LocalEarlyLoadedSettingSource(fname, dir)
	v, err = ss.Get(ctx, "c1")
	assert.Equal(t, "ValueFromFile", tricks.PtrVal(v))
	require.NoError(t, err)

	cfgTwo := struct {
		C2 string `json:"c2"`
	}{
		C2: "AnotheValueInFile",
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
}
