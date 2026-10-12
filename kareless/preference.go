package kareless

import "context"

// Preferences stores caller-owned string settings that outlive a single request.
type Preferences interface {
	// SetPreference stores value at key.
	SetPreference(ctx context.Context, key, value string) error
	// GetPreference returns the value stored at key.
	GetPreference(ctx context.Context, key string) (string, error)
}
