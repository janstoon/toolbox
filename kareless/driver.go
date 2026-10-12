package kareless

import "context"

type (
	// Driver runs an installed application, for example an HTTP server or a queue consumer.
	Driver interface {
		// Run serves until ctx is canceled or the driver fails.
		Run(ctx context.Context) error

		// todo: Start(ctx context.Context) error
		// todo: Shutdown(ctx context.Context) error
	}
	// DriverConstructor builds one Driver from settings, instruments, and the installed applications.
	DriverConstructor func(ss *Settings, ib *InstrumentBank, apps []Application) Driver
)
