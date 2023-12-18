# Janstun Toolbox [![Build Status][bdg-build-status]][action-tests]

A bunch of converters, data structures, parsers, validators, middlewares and injectors useful when crafting software in
Golang specially in early stages. They're designed with **Open-closed principle** in mind; In other words they are open
for extension and customization by others from outside. Most of them introduce and work with interfaces or data types
which are implemented for specific use-cases while can become implemented and extended in external packages.

This repository is a monorepo which consists of some (go) modules, each providing related functionalities:

## [Bricks][mod-bricks] [![Coverage Badge][bdg-cov-bricks]][action-tests]
Models of real world basic entities like phone number and iban.

## [Tricks][mod-tricks] [![Coverage Badge][bdg-cov-tricks]][action-tests]
Frequent low-level functionalities like conversions or functional paradigm missing utilities.

## [Handywares][mod-handywares] [![Coverage Badge][bdg-cov-handywares]][action-tests]
Mid-level requirements designed in a reusable manner like general http middlewares.

## [Kareless][mod-kareless] [![Coverage Badge][bdg-cov-kareless]][action-tests]
Thin dependency injector which tries not to be a dependency injector at all.

# License
This library is [licensed](LICENSE) under the [GPL v3 License][gpl]. © 2023 [Janstun][janstun]

[action-tests]: https://github.com/janstoon/toolbox/actions?query=branch%3Amaster+workflow%3Atests
[bdg-build-status]: https://github.com/janstoon/toolbox/actions/workflows/tests.yml/badge.svg?branch=master
[bdg-cov-tricks]: https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/pouyanh/69229998008a13b9b87590ebe50ecded/raw/janstoon_toolbox_tricks_refs_heads_master.json
[bdg-cov-bricks]: https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/pouyanh/69229998008a13b9b87590ebe50ecded/raw/janstoon_toolbox_bricks_refs_heads_master.json
[bdg-cov-handywares]: https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/pouyanh/69229998008a13b9b87590ebe50ecded/raw/janstoon_toolbox_handywares_refs_heads_master.json
[bdg-cov-kareless]: https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/pouyanh/69229998008a13b9b87590ebe50ecded/raw/janstoon_toolbox_kareless_refs_heads_master.json
[mod-bricks]: bricks
[mod-tricks]: tricks
[mod-handywares]: handywares
[mod-kareless]: kareless
[gpl]: https://www.gnu.org/licenses/gpl-3.0.en.html
[janstun]: http://janstun.com
