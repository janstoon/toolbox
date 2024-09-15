package kareless

type (
	Application interface {
		// todo: a method to let in-progress jobs to finish and wait for them
		// todo: optional reload method to signal applications to reload
	}
	ApplicationConstructor func(ss *Settings, ib *InstrumentBank) Application
)
