// Package infra generates eda_software_infrastructure problems: design databases and object
// models, tool APIs and DSLs, format parsers, tool architecture, incremental computation,
// parallelism and runs, testing and builds, open PDKs, the open-source flow and the agent-first
// frontier of EDA tool construction. Files are split by subdomain; each file adds its own
// register function to registrars from an init function.
package infra

import "github.com/chuma-beep/mathua/internal/generator"

// registrars holds one register function per subdomain file.
var registrars []func(*generator.Registry)

// Register adds every eda_software_infrastructure generator to reg.
func Register(reg *generator.Registry) {
	for _, r := range registrars {
		r(reg)
	}
}
