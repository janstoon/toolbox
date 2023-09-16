package kareless

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Kernel struct {
	ss *Settings
	ib *InstrumentBank

	appsToInstall    []ApplicationConstructor
	driversToConnect []DriverConstructor
	postHooks        []Hook
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

func (k Kernel) Run(ctx context.Context) error {
	apps := make([]Application, len(k.appsToInstall))
	for i, constructor := range k.appsToInstall {
		apps[i] = constructor(k.ss, k.ib)
	}

	wgDrivers, wgAll := sync.WaitGroup{}, errgroup.Group{}
	for _, constructor := range k.driversToConnect {
		wgDrivers.Add(1)
		func(constructor DriverConstructor) {
			wgAll.Go(func() error {
				driver := constructor(k.ss, k.ib, apps)
				wgDrivers.Done()

				return driver.Run(ctx)
			})
		}(constructor)
	}

	wgDrivers.Wait()
	for _, hook := range k.postHooks {
		func(hook Hook) {
			wgAll.Go(func() error {
				return hook(ctx, k.ss, k.ib, apps)
			})
		}(hook)
	}

	return wgAll.Wait()
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
	k.appsToInstall = append(k.appsToInstall, cc...)

	return k
}

func Connector(cc ...DriverConstructor) Option {
	return func(k Kernel) Kernel {
		return k.Connect(cc...)
	}
}

func (k Kernel) Connect(cc ...DriverConstructor) Kernel {
	k.driversToConnect = append(k.driversToConnect, cc...)

	return k
}

type Hook func(ctx context.Context, ss *Settings, ib *InstrumentBank, apps []Application) error

func PostHook(hh ...Hook) Option {
	return func(k Kernel) Kernel {
		return k.AfterStart(hh...)
	}
}

func (k Kernel) AfterStart(hh ...Hook) Kernel {
	k.postHooks = append(k.postHooks, hh...)

	return k
}
