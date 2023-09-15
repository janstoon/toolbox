package kareless

import "context"

type (
	Driver interface {
		Start(ctx context.Context) error
	}
	DriverConstructor func(ss *Settings, ib *InstrumentBank, apps []Application) Driver
)
