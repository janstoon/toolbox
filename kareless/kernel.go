package kareless

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

// Run creates installed applications and connected drivers and waits until drivers and hooks are all finished running
func (k Kernel) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancelCause(ctx)
	bindSignals(func(sig os.Signal) {
		// todo: Do not cancel the driver.Run's context as a shutdown signal.
		// Because drivers have to be able to continue finishing requests procedures to perform a graceful shutdown
		// Instead call driver.Stop concurrently.
		// Shutdown sequence:
		//   1. Drivers: to eliminate new requests acceptance.
		//   2. Applications: to finish in-progress jobs
		cancel(fmt.Errorf("signal caught: %s. context canceled", sig))
	}, syscall.SIGTERM, syscall.SIGINT /* todo: Reload on interrupt */)

	if err := k.ib.openCatalogues(k.ss); err != nil {
		return err
	}

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

func bindSignals(fn func(sig os.Signal), ss ...os.Signal) {
	ntfy := make(chan os.Signal, 1)
	signal.Notify(ntfy, ss...)

	go func() {
		fn(<-ntfy)
	}()
}

func Feeder(ss ...SettingSource) Option {
	return func(k Kernel) Kernel {
		return k.Feed(ss...)
	}
}

// Feed injects setting sources to be fed into configurable units like instruments, applications and drivers
func (k Kernel) Feed(ss ...SettingSource) Kernel {
	for _, source := range ss {
		k.ss.Append(source)
	}

	return k
}

func Equipment(cc ...InstrumentInjector) Option {
	return func(k Kernel) Kernel {
		return k.Equip(cc...)
	}
}

// Equip plugs instruments which can get resolved by the instrument bank that is passed to unit constructors
func (k Kernel) Equip(cc ...InstrumentInjector) Kernel {
	k.ib.register(cc...)

	return k
}

func Installer(cc ...ApplicationConstructor) Option {
	return func(k Kernel) Kernel {
		return k.Install(cc...)
	}
}

// Install appends installable applications to the list which become created on Run
func (k Kernel) Install(cc ...ApplicationConstructor) Kernel {
	k.appsToInstall = append(k.appsToInstall, cc...)

	return k
}

func Connector(cc ...DriverConstructor) Option {
	return func(k Kernel) Kernel {
		return k.Connect(cc...)
	}
}

// Connect binds driver(s) to the Kernel in order to invoke use-cases on (drive) installed applications
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
