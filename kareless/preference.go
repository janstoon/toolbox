package kareless

import "context"

type Preferences interface {
	SetPreference(ctx context.Context, key, value string) error
	GetPreference(ctx context.Context, key string) (string, error)
}
