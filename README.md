# Toolbox [![Build Status][bdg-build-status]][action-tests]

A bunch of converters, data structures, parsers, validators, middlewares and injectors useful when crafting software in
Golang specially in early stages. They're designed with **Open-closed principle** in mind; In other words they are open
for extension and customization by others from outside. Most of them introduce and work with interfaces or data types
which are implemented for specific use-cases while can become implemented and extended in external packages.

This repository consists of some (go) modules, each providing related functionalities:
## Bricks [![Coverage Badge][bdg-cov-bricks]][action-tests]
Models of real world basic entities like phone number and iban.
## Tricks [![Coverage Badge][bdg-cov-tricks]][action-tests]
Frequent low-level functionalities like conversions or functional paradigm missing utilities.
## Handywares [![Coverage Badge][bdg-cov-handywares]][action-tests]
Mid-level requirements designed in a reusable manner like general http middlewares.
## Kareless [![Coverage Badge][bdg-cov-kareless]][action-tests]
Thin dependency injector which tries not to be a dependency injector at all.

[action-tests]: https://github.com/janstoon/toolbox/actions?query=branch%3Amaster+workflow%3Atests
[bdg-build-status]: https://github.com/janstoon/toolbox/actions/workflows/tests.yml/badge.svg?branch=master
[bdg-cov-tricks]: https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/pouyanh/69229998008a13b9b87590ebe50ecded/raw/janstoon_toolbox_tricks_refs_heads_master.json
[bdg-cov-bricks]: https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/pouyanh/69229998008a13b9b87590ebe50ecded/raw/janstoon_toolbox_bricks_refs_heads_master.json
[bdg-cov-handywares]: https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/pouyanh/69229998008a13b9b87590ebe50ecded/raw/janstoon_toolbox_handywares_refs_heads_master.json
[bdg-cov-kareless]: https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/pouyanh/69229998008a13b9b87590ebe50ecded/raw/janstoon_toolbox_kareless_refs_heads_master.json
