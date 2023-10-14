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

	"github.com/janstoon/toolbox/kareless/std"
)

func TestLocalSettings_Get(t *testing.T) {
	dir, err := os.MkdirTemp("", "kareless-std-local-settings-")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(dir) }()

	fname := "tconf"
	fh, err := os.Create(path.Join(dir, fmt.Sprintf("%s.json", fname)))
	assert.NoError(t, err)
	defer func() { _ = fh.Close() }()

	enc := json.NewEncoder(fh)
	ctx := context.Background()

	ss := std.LocalEarlyLoadedSettingSource("")
	v, err := ss.Get(ctx, "c1")
	assert.Empty(t, v)
	assert.ErrorIs(t, err, bricks.ErrNotFound)

	cfgOne := struct {
		C1 string `json:"c1"`
	}{
		C1: "ValueFromFile",
	}
	err = enc.Encode(cfgOne)
	assert.NoError(t, err)
	err = fh.Sync()
	assert.NoError(t, err)
	ss = std.LocalEarlyLoadedSettingSource(fname, dir)
	v, err = ss.Get(ctx, "c1")
	assert.Equal(t, "ValueFromFile", tricks.PtrVal(v))
	assert.NoError(t, err)

	err = os.Setenv("C1", "ValueFromEnv")
	assert.NoError(t, err)
	ss = std.LocalEarlyLoadedSettingSource(fname, dir)
	v, err = ss.Get(ctx, "c1")
	assert.Equal(t, "ValueFromEnv", tricks.PtrVal(v))
	assert.NoError(t, err)

	err = os.Unsetenv("C1")
	assert.NoError(t, err)
	ss = std.LocalEarlyLoadedSettingSource(fname, dir)
	v, err = ss.Get(ctx, "c1")
	assert.Equal(t, "ValueFromFile", tricks.PtrVal(v))
	assert.NoError(t, err)

	cfgTwo := struct {
		C2 string `json:"c2"`
	}{
		C2: "AnotheValueInFile",
	}
	err = fh.Truncate(0)
	assert.NoError(t, err)
	_, err = fh.Seek(0, 0)
	assert.NoError(t, err)
	err = enc.Encode(cfgTwo)
	assert.NoError(t, err)
	err = fh.Sync()
	assert.NoError(t, err)
	ss = std.LocalEarlyLoadedSettingSource(fname, dir)
	v, err = ss.Get(ctx, "c1")
	assert.Empty(t, v)
	assert.ErrorIs(t, err, bricks.ErrNotFound)
}
