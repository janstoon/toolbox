package kareless

type (
	// Application is a use-case the Kernel installs and passes to drivers.
	// It has no required methods so each application can expose its own API.
	Application interface {
		// todo: a method to let in-progress jobs to finish and wait for them
		// todo: optional reload method to signal applications to reload
	}
	// ApplicationConstructor builds one Application from settings and the instrument bank.
	ApplicationConstructor func(ss *Settings, ib *InstrumentBank) Application
)
