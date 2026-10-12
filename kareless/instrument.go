package kareless

import (
	"errors"
	"fmt"
	"sync"
)

var (
	// ErrAlreadyRegisteredInstrument is returned when two catalogues use the same name.
	ErrAlreadyRegisteredInstrument = errors.New("instrument already registered")
	// ErrUnresolvedDependency is the panic value when Resolve is asked for a missing name.
	ErrUnresolvedDependency = errors.New("dependency not resolved")
	// ErrUnacceptableDependency is the panic value when the tester rejects the resolved value.
	ErrUnacceptableDependency = errors.New("dependency not acceptable")
)

type (
	// InstrumentInjector lists the instruments one module contributes. Equip calls it while opening catalogues.
	InstrumentInjector func(ss *Settings) []InstrumentCatalogue
	// InstrumentCatalogue publishes one instrument under one or more names.
	InstrumentCatalogue struct {
		// Names are the lookup keys. Each name must be unique across the bank.
		Names []string
		// Builder constructs the instrument the first time a name is resolved.
		Builder InstrumentConstructor
	}
	// InstrumentConstructor builds an instrument. The result is cached for every name in the catalogue.
	InstrumentConstructor func(ss *Settings, ib *InstrumentBank) Instrument
	// Instrument is any dependency a constructor returns. Callers assert the concrete type.
	Instrument interface{}
)

type instrumentFactory struct {
	sync.Once
	constructor InstrumentConstructor
	cached      Instrument
}

func (fkt *instrumentFactory) create(ss *Settings, ib *InstrumentBank) Instrument {
	fkt.Do(func() {
		// recover possible panic

		fkt.cached = fkt.constructor(ss, ib)
	})

	return fkt.cached
}

// InstrumentBank resolves named dependencies. Each name is built at most once.
type InstrumentBank struct {
	lock      sync.RWMutex
	injectors []InstrumentInjector
	factories map[string]*instrumentFactory

	settings *Settings
}

// Resolve returns the instrument stored under name.
// tester reports whether the concrete value is acceptable.
// Resolve panics with ErrUnresolvedDependency when name is missing and with
// ErrUnacceptableDependency when tester returns false.
func (ib *InstrumentBank) Resolve(name string, tester func(v any) bool) Instrument {
	ins := ib.resolve(name)
	if ins == nil {
		panic(errors.Join(ErrUnresolvedDependency, fmt.Errorf("no instrument provided: %s", name)))
	}

	if !tester(ins) {
		panic(errors.Join(ErrUnacceptableDependency, fmt.Errorf("test failed for instrument(%s): %T", name, ins)))
	}

	return ins
}

func newInstrumentBank(ss *Settings) *InstrumentBank {
	return &InstrumentBank{
		factories: make(map[string]*instrumentFactory),
		settings:  ss,
	}
}

func (ib *InstrumentBank) register(cc ...InstrumentInjector) {
	ib.lock.Lock()
	defer ib.lock.Unlock()

	ib.injectors = append(ib.injectors, cc...)
}

func (ib *InstrumentBank) openCatalogues(ss *Settings) error {
	ib.lock.Lock()
	defer ib.lock.Unlock()

	ib.factories = make(map[string]*instrumentFactory)
	for _, injector := range ib.injectors {
		catalogues := injector(ss)
		for _, catalogue := range catalogues {
			fkt := &instrumentFactory{
				constructor: catalogue.Builder,
			}

			for _, name := range catalogue.Names {
				if _, registered := ib.factories[name]; registered {
					return errors.Join(ErrAlreadyRegisteredInstrument, fmt.Errorf("duplicate entry for: %s", name))
				}

				ib.factories[name] = fkt
			}
		}
	}

	return nil
}

func (ib *InstrumentBank) resolve(name string) Instrument {
	ib.lock.RLock()
	fkt, ok := ib.factories[name]
	ib.lock.RUnlock()
	if !ok || fkt == nil {
		return nil
	}

	return fkt.create(ib.settings, ib)
}

// ResolveInstrumentByType returns the instrument named name asserted to T.
// It panics when the name is missing or the value is not a T.
func ResolveInstrumentByType[T any](ib *InstrumentBank, name string) T {
	i, _ := ib.Resolve(name, InstrumentTesterByTypeAssertion[T]).(T)

	return i
}

// InstrumentTesterByTypeAssertion reports whether v has type T.
func InstrumentTesterByTypeAssertion[T any](v any) bool {
	_, ok := v.(T)

	return ok
}
