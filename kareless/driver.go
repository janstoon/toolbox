package kareless

import "context"

type (
	Driver interface {
		Run(ctx context.Context) error

		// todo: Start(ctx context.Context) error
		// todo: Shutdown(ctx context.Context) error
	}
	DriverConstructor func(ss *Settings, ib *InstrumentBank, apps []Application) Driver
)
