package std

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/janstoon/toolbox/bricks"
	"github.com/spf13/viper"

	"github.com/janstoon/toolbox/kareless"
)

type MapSettingSource map[string]any

func (ss MapSettingSource) Get(ctx context.Context, key string) (any, error) {
	if kareless.IsSettingRootKey(key) {
		return ss, nil
	}

	v, foundDirect := ss[key]
	if foundDirect {
		return v, nil
	} else if path := strings.SplitN(key, kareless.SettingKeyDelimiter, 2); len(path) > 1 {
		if parent, err := ss.navigable(path[0]); err == nil {
			return parent.Get(ctx, path[1])
		}
	}

	return nil, bricks.ErrNotFound
}

func (ss MapSettingSource) navigable(root string) (MapSettingSource, error) {
	parent, foundParent := ss[root]
	if !foundParent {
		return nil, bricks.ErrNotFound
	}

	var mss MapSettingSource
	bb, err := json.Marshal(parent)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bb, &mss)
	if err != nil {
		return nil, err
	}

	return mss, nil
}

type localEarlySettings struct {
	v *viper.Viper
}

// LocalEarlyLoadedSettingSource provides a setting source fed from files. Keys are case-insensitive.
// paths are directories and name is filename without extension. Files can be in any supported formats including
// json, yaml and toml with appropriate extension (.json, .yml, .yaml, .toml).
func LocalEarlyLoadedSettingSource(name string, paths ...string) kareless.SettingSource {
	v := viper.NewWithOptions(viper.EnvKeyReplacer(strings.NewReplacer(kareless.SettingKeyDelimiter, "_")))
	v.AutomaticEnv()
	v.SetConfigName(name)
	for _, p := range paths {
		v.AddConfigPath(p)
	}

	if len(paths) > 0 {
		if err := v.ReadInConfig(); err != nil {
			if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
				panic(err)
			}
		}
	}

	return localEarlySettings{
		v: v,
	}
}

func (ss localEarlySettings) Get(_ context.Context, key string) (any, error) {
	if kareless.IsSettingRootKey(key) {
		return ss.v.AllSettings(), nil
	}

	if ss.v.IsSet(key) {
		return ss.v.Get(key), nil
	}

	return nil, bricks.ErrNotFound
}
