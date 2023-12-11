# Toolbox
A bunch of converters, data structures, parsers, validators, middlewares and injectors useful when crafting software in
Golang specially in early stages. They're designed with **Open-closed principle** in mind; In other words they are open
for extension and customization by others from outside. Most of them introduce and work with interfaces or data types
which are implemented for specific use-cases while can become implemented and extended in external packages.

This repository consists of some (go) modules, each providing related functionalities:
* Tricks: Frequent low-level functionalities like conversions or functional paradigm missing utilities.
* Bricks: Models of real world basic entities like phone number and iban.
* Handywares: Mid-level requirements designed in a reusable manner like general http middlewares.
* Kareless: Thin dependency injector which tries not to be a dependency injector at all.
