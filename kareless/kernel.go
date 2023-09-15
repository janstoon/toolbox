package kareless

import "context"

type Kernel struct {
	ss   *Settings
	ib   *InstrumentBank
	apps []Application
}

type Option func(k Kernel) Kernel

func Compile(oo ...Option) Kernel {
	k := Kernel{
		ss: new(Settings),
	}
	k.ib = newInstrumentBank(k.ss)

	for _, o := range oo {
		k = o(k)
	}

	return k
}

func (k Kernel) Run(ctx context.Context) {

}

func Feeder(ss ...SettingSource) Option {
	return func(k Kernel) Kernel {
		return k.Feed(ss...)
	}
}

func (k Kernel) Feed(ss ...SettingSource) Kernel {
	for _, source := range ss {
		k.ss.append(source)
	}

	return k
}

func Equipment(cc ...InstrumentCatalogue) Option {
	return func(k Kernel) Kernel {
		return k.Equip(cc...)
	}
}

func (k Kernel) Equip(cc ...InstrumentCatalogue) Kernel {
	for _, catalogue := range cc {
		k.ib.register(catalogue)
	}

	return k
}

func Installer(cc ...ApplicationConstructor) Option {
	return func(k Kernel) Kernel {
		return k.Install(cc...)
	}
}

func (k Kernel) Install(cc ...ApplicationConstructor) Kernel {
	for _, constructor := range cc {
		k.apps = append(k.apps, constructor(k.ss, k.ib))
	}

	return k
}

func Connector(cc ...DriverConstructor) Option {
	return func(k Kernel) Kernel {
		return k.Connect(cc...)
	}
}

func (k Kernel) Connect(cc ...DriverConstructor) Kernel {
	for _, constructor := range cc {
		constructor(k.ss, k.ib, k.apps)
	}

	return k
}
