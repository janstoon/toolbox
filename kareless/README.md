# Kareless [![Coverage Badge][bdg-cov-kareless]][action-tests]
A pico-framework to glue building blocks of a long-running software together.

## Why another dependency injector at all?
These are the reasons why I need another _di_, while there are great solutions like [uber's fx][uber-fx] in the wild:
* Dependant (i.e. consumer) defines the expected type rather than the dependency itself;
In other words, the dependant tries to cast the available dependency, which is resolved by name, to the desirable type
and can use it in case of success.
* Multi-tier settings storage is a first-class citizen.

# Usage

```shell
go get github.com/janstoon/toolbox/kareless
```

[action-tests]: https://github.com/janstoon/toolbox/actions?query=branch%3Amaster+workflow%3Atests
[bdg-cov-kareless]: https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/pouyanh/69229998008a13b9b87590ebe50ecded/raw/janstoon_toolbox_kareless_refs_heads_master.json
[uber-fx]: https://go.uber.org/fx
