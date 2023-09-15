package kareless

type (
	Driver            interface{}
	DriverConstructor func(ss *Settings, ib *InstrumentBank, apps []Application) Driver
)
