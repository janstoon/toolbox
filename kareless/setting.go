package kareless

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cast"
)

// SettingSource returns the value stored at key.
// A missing key returns an error. A nil value is treated as missing by Settings.
type SettingSource interface {
	// Get returns the value stored at key.
	Get(ctx context.Context, key string) (any, error)
}

// SettingKeyDelimiter separates nested setting keys, as in "redis.address".
const SettingKeyDelimiter = "."

// IsSettingRootKey reports whether key names the whole source rather than one entry.
// An empty key and "." are root keys.
func IsSettingRootKey(key string) bool {
	return len(strings.TrimSpace(key)) == 0 || key == "."
}

// Settings reads the first source that returns a non-nil value.
// Prepend puts a source ahead of the others. Append puts it behind them.
type Settings struct {
	lock sync.RWMutex
	rr   []SettingSource
}

// Prepend inserts source at the front, so it is consulted before sources added earlier.
func (ss *Settings) Prepend(source SettingSource) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	ss.rr = append([]SettingSource{source}, ss.rr...)
}

// Append inserts source at the end, so it is a fallback for sources added earlier.
func (ss *Settings) Append(source SettingSource) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	ss.rr = append(ss.rr, source)
}

// UnmarshalJson JSON-encodes the value at key and decodes it into valPtr.
func (ss *Settings) UnmarshalJson(key string, valPtr any) error {
	bb, err := json.Marshal(ss.get(context.Background(), key))
	if err != nil {
		return err
	}

	return json.Unmarshal(bb, valPtr)
}

// GetString returns the value at key as a string. A missing key returns the empty string.
func (ss *Settings) GetString(key string) string {
	return cast.ToString(ss.get(context.Background(), key))
}

// GetStringSlice returns the value at key as a slice of strings. A missing key returns nil.
func (ss *Settings) GetStringSlice(key string) []string {
	return cast.ToStringSlice(ss.get(context.Background(), key))
}

// GetInt returns the value at key as an int. A missing or unparsable key returns 0.
func (ss *Settings) GetInt(key string) int {
	return cast.ToInt(ss.get(context.Background(), key))
}

// GetInt64 returns the value at key as an int64. A missing or unparsable key returns 0.
func (ss *Settings) GetInt64(key string) int64 {
	return cast.ToInt64(ss.get(context.Background(), key))
}

// GetByte returns the value at key as a byte. A missing or unparsable key returns 0.
func (ss *Settings) GetByte(key string) byte {
	return byte(ss.GetInt(key))
}

// GetBool returns the value at key as a bool. A missing key returns false.
func (ss *Settings) GetBool(key string) bool {
	return cast.ToBool(ss.get(context.Background(), key))
}

// GetDuration returns the value at key as a time.Duration. A missing key returns 0.
func (ss *Settings) GetDuration(key string) time.Duration {
	return cast.ToDuration(ss.get(context.Background(), key))
}

// Children returns the indexes of a slice value or the keys of a map value stored at key.
// It returns nil when the value is missing or is neither a slice nor a map.
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
