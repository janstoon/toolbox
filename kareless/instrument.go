package kareless

import "sync"

type (
	Instrument            interface{}
	InstrumentConstructor func(ss *Settings, ib *InstrumentBank) Instrument
	InstrumentCatalogue   func() ([]string, InstrumentConstructor)
)

type instrumentFactory struct {
	sync.Once
	constructor InstrumentConstructor
	cached      Instrument
}

func (fkt *instrumentFactory) create(ss *Settings, ib *InstrumentBank) Instrument {
	fkt.Do(func() {
		fkt.cached = fkt.constructor(ss, ib)
	})

	return fkt.cached
}

type InstrumentBank struct {
	lock      sync.RWMutex
	factories map[string]*instrumentFactory

	settings *Settings
}

func newInstrumentBank(ss *Settings) *InstrumentBank {
	return &InstrumentBank{
		factories: make(map[string]*instrumentFactory),
		settings:  ss,
	}
}

func (ib *InstrumentBank) register(catalogue InstrumentCatalogue) {
	names, ic := catalogue()
	fkt := &instrumentFactory{
		constructor: ic,
	}

	ib.lock.Lock()
	defer ib.lock.Unlock()
	for _, name := range names {
		if _, registered := ib.factories[name]; registered {
			panic("instrument already registered")
		}

		ib.factories[name] = fkt
	}
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
		panic("dependency not resolved")
	}

	if !tester(ins) {
		panic("dependency unacceptable")
	}

	return ins
}