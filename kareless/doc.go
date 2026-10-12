// Package kareless is a small runtime that wires settings, named instruments,
// applications, and drivers without reflection.
//
// Assemble a Kernel with Compile, then Feed setting sources, Equip instruments,
// Install applications, and Connect drivers. Run starts drivers and after-start hooks.
// Dependencies are stored by name and retrieved with a type assertion.
package kareless
