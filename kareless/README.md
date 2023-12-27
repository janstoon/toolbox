# Kareless [![Coverage Badge][bdg-cov-kareless]][action-tests]
A pico-framework to glue building blocks of a long-running software together.

## Why another dependency injector at all?
These are the reasons why I need another _di_, while there are great solutions like [uber's fx][uber-fx] in the wild:
* Dependant (i.e. consumer) defines the expected type rather than the dependency itself;
In other words, the dependant tries to cast the available dependency, which is resolved by name, to the desirable type
and can use it in case of success.
* Multi-tier settings storage is a first-class citizen.
* Put restriction on cross-layer dependency to prevent making a mesh.
* Reflection-free dependency resolution using string names and typecasts
* Flexible dependency graph using settings to reduce unnecessary code changes and rebuilds

# Usage

```shell
go get github.com/janstoon/toolbox/kareless
```

```go
package main

import (
	"context"
	"log"

	"github.com/janstoon/toolbox/kareless"
	"github.com/janstoon/toolbox/kareless/std"
	"github.com/redis/go-redis/v9"
)

func main() {
	k := kareless.Compile().
		Feed( /* L1 Settings Source */).
		Feed(std.LocalEarlyLoadedSettingSource("conf", "/etc/janstoon")).
		Feed( /* L2 Settings Source */).
		Feed(std.MapSettingSource{
			"log.level":     5,
			"debug":         true,
			"redis.address": "redis://redis.janstun.com:6379/1",
		}).
		Feed( /* Default Settings Source as a Fallback */).
		Equip( /* Instruments as dependencies resolvable by name */).
		Equip(func(ss *kareless.Settings, ib *kareless.InstrumentBank) []kareless.InstrumentCatalogue {
			return []kareless.InstrumentCatalogue{
				{
					Names: []string{"redis", "infra/redis"},
					Builder: func(ss *kareless.Settings, ib *kareless.InstrumentBank) kareless.Instrument {
						oo, err := redis.ParseURL(ss.GetString("redis.address"))
						if err != nil {
							panic(err)
						}

						return redis.NewClient(oo)
					},
				},
			}
		}).
		Equip( /* Instruments as dependencies resolvable by name */).
		Install( /* Applications as dependencies available in a slice */).
		Install(
			func(ss *kareless.Settings, ib *kareless.InstrumentBank) kareless.Application {
				return appFoo{
					rc: kareless.ResolveInstrumentByType[*redis.Client](ib, "redis"),
				}
			},
			func(ss *kareless.Settings, ib *kareless.InstrumentBank) kareless.Application {
				return appBar{
					rc: kareless.ResolveInstrumentByType[*redis.Client](ib, "infra/redis"),
				}
			},
		).
		Connect( /* Drivers to drive application(s) */).
		Connect(func(ss *kareless.Settings, ib *kareless.InstrumentBank, apps []kareless.Application) kareless.Driver {
			type (
				appAccess = appFoo
				appOrdering  = appBar
			)

			adapter := struct {
				appAccess
				appOrdering
			}{}

			for _, v := range apps {
				switch app := v.(type) {
				case appFoo:
					adapter.appAccess = app

				case appBar:
					adapter.appOrdering = app
				}
			}

			return gateway{
				app: adapter,
			}
		}).
		AfterStart( /* Hooks to run after drivers started */).
		AfterStart(
			func(ctx context.Context, ss *kareless.Settings, ib *kareless.InstrumentBank, apps []kareless.Application) error {
				log.Println("System started...")

				return nil
			},
		)
	if err := k.Run(context.Background()); err != nil {
		panic(err)
	}
}

// ----------------------
//     Application(s)
// ----------------------

type appFoo struct {
	rc *redis.Client
}

func (a appFoo) Authenticate(ctx context.Context, token string) error {
	// Handle the use-case

	return nil
}

type appBar struct {
	rc *redis.Client
}

func (a appBar) SubmitOrder(ctx context.Context, details any) error {
	// Handle the use-case

	return nil
}

// ----------------------
//       Driver(s)
// ----------------------

type gwApp interface {
	accessApp
	orderingApp
}

type accessApp interface {
	Authenticate(ctx context.Context, token string) error
}

type orderingApp interface {
	SubmitOrder(ctx context.Context, details any) error
}

type gateway struct {
	app gwApp
}

func (d gateway) Run(ctx context.Context) error {
	// Start serving the service and driving the application

	return nil
}
```

# Concepts

## Setting Source
## Instrument
## Driver
## Application
## Kernel

[action-tests]: https://github.com/janstoon/toolbox/actions?query=branch%3Amaster+workflow%3Atests
[bdg-cov-kareless]: https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/pouyanh/69229998008a13b9b87590ebe50ecded/raw/janstoon_toolbox_kareless_refs_heads_master.json
[uber-fx]: https://go.uber.org/fx
