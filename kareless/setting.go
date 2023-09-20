package kareless

import (
	"context"
	"strconv"
	"sync"

	"github.com/janstoon/toolbox/conv"
)

type SettingSource interface {
	Get(ctx context.Context, key string) (*string, error)
}

type Settings struct {
	lock sync.RWMutex
	rr   []SettingSource
}

func (ss *Settings) append(source SettingSource) {
	ss.lock.Lock()
	defer ss.lock.Unlock()
	ss.rr = append(ss.rr, source)
}

func (ss *Settings) get(ctx context.Context, key string) string {
	ss.lock.RLock()
	defer ss.lock.RUnlock()

	for _, r := range ss.rr {
		v, err := r.Get(ctx, key)
		if err == nil && v != nil {
			return conv.PtrVal(v)
		}
	}

	return ""
}

func (ss *Settings) GetString(key string) string {
	return ss.get(context.Background(), key)
}

func (ss *Settings) GetInt(key string) int {
	v, err := strconv.Atoi(ss.GetString(key))
	if err != nil {
		return 0
	}

	return v
}

func (ss *Settings) GetInt64(key string) int64 {
	return int64(ss.GetInt(key))
}

func (ss *Settings) GetByte(key string) byte {
	return byte(ss.GetInt(key))
}
