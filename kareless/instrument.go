package kareless

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrAlreadyRegisteredInstrument = errors.New("instrument already registered")
	ErrUnresolvedDependency        = errors.New("dependency not resolved")
	ErrUnacceptableDependency      = errors.New("dependency not acceptable")
)

type (
	Instrument            interface{}
	InstrumentConstructor func(ss *Settings, ib *InstrumentBank) Instrument
	InstrumentCatalogue   func(ss *Settings, ib *InstrumentBank) ([]string, InstrumentConstructor)
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

type InstrumentBank struct {
	lock       sync.RWMutex
	catalogues []InstrumentCatalogue
	factories  map[string]*instrumentFactory

	settings *Settings
}

func newInstrumentBank(ss *Settings) *InstrumentBank {
	return &InstrumentBank{
		factories: make(map[string]*instrumentFactory),
		settings:  ss,
	}
}

func (ib *InstrumentBank) register(cc ...InstrumentCatalogue) {
	ib.lock.Lock()
	defer ib.lock.Unlock()

	ib.catalogues = append(ib.catalogues, cc...)
}

func (ib *InstrumentBank) openCatalogues(ss *Settings) error {
	ib.lock.Lock()
	defer ib.lock.Unlock()

	for _, catalogue := range ib.catalogues {
		names, ic := catalogue(ss, ib)
		fkt := &instrumentFactory{
			constructor: ic,
		}

		for _, name := range names {
			if _, registered := ib.factories[name]; registered {
				return errors.Join(ErrAlreadyRegisteredInstrument, fmt.Errorf("duplicate entry for: %s", name))
			}

			ib.factories[name] = fkt
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

func (ib *InstrumentBank) Resolve(name string, tester func(v any) bool) Instrument {
	ins := ib.resolve(name)
	if ins == nil {
		panic(errors.Join(ErrUnresolvedDependency, fmt.Errorf("no instrument provided: %s", name)))
	}

	if !tester(ins) {
		panic(errors.Join(ErrUnacceptableDependency, fmt.Errorf("test failed dor instrument(%s): %T", name, ins)))
	}

	return ins
}

func ResolveInstrumentByType[T any](ib *InstrumentBank, name string) T {
	i, _ := ib.Resolve(name, InstrumentTesterByTypeAssertion[T]).(T)

	return i
}

func InstrumentTesterByTypeAssertion[T any](v any) bool {
	_, ok := v.(T)

	return ok
}
