package kareless

import (
	"context"
	"sync"
	"time"

	"github.com/spf13/cast"
)

type SettingSource interface {
	Get(ctx context.Context, key string) (any, error)
}

type Settings struct {
	lock sync.RWMutex
	rr   []SettingSource
}

func (ss *Settings) Prepend(source SettingSource) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	ss.rr = append([]SettingSource{source}, ss.rr...)
}

func (ss *Settings) Append(source SettingSource) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	ss.rr = append(ss.rr, source)
}

func (ss *Settings) get(ctx context.Context, key string) any {
	ss.lock.RLock()
	defer ss.lock.RUnlock()

	for _, r := range ss.rr {
		v, err := r.Get(ctx, key)
		if err == nil && v != nil {
			return v
		}
	}

	return nil
}

func (ss *Settings) GetString(key string) string {
	return cast.ToString(ss.get(context.Background(), key))
}

func (ss *Settings) GetStringSlice(key string) []string {
	return cast.ToStringSlice(ss.get(context.Background(), key))
}

func (ss *Settings) GetInt(key string) int {
	return cast.ToInt(ss.get(context.Background(), key))
}

func (ss *Settings) GetInt64(key string) int64 {
	return cast.ToInt64(ss.get(context.Background(), key))
}

func (ss *Settings) GetByte(key string) byte {
	return byte(ss.GetInt(key))
}

func (ss *Settings) GetBool(key string) bool {
	return cast.ToBool(ss.get(context.Background(), key))
}

func (ss *Settings) GetDuration(key string) time.Duration {
	return cast.ToDuration(ss.get(context.Background(), key))
}

func (ss *Settings) Children(key string) []string {
	v := ss.get(context.Background(), key)
	if aa, err := cast.ToSliceE(v); err == nil {
		kk := make([]string, len(aa))
		for k := range aa {
			kk[k] = cast.ToString(k)
		}

		return kk
	}

	if aa, err := cast.ToStringMapE(v); err == nil {
		kk := make([]string, 0, len(aa))
		for k := range aa {
			kk = append(kk, k)
		}

		return kk
	}

	return nil
}
