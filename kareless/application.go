package kareless

type (
	Application            interface{}
	ApplicationConstructor func(ss *Settings, ib *InstrumentBank) Application
)
