# UPF power intent

The power architecture of a chip, its domains, supplies, switches, isolation, level shifting and retention, is not written in the RTL. It is captured separately as power intent in the Unified Power Format (UPF, IEEE 1801), or in the older Common Power Format (CPF). Every tool in the flow reads the same intent: power-aware simulation, synthesis, place and route, and signoff. A tool designer must know what each construct means, because the tool must implement and verify it.

## Domains and supplies

create_power_domain groups instances that share a primary supply and switch together. create_supply_port and create_supply_net declare supplies, connect_supply_net wires them, and the domain's primary power and ground are attached to it. create_power_switch declares a header or footer switch with its input supply, output (switched) supply and control port, which implementation turns into real switch cells. Keeping intent outside the RTL means the same logic can be verified and implemented under different power architectures. CPF, from the Si2 consortium, expresses the same concepts; most flows now use UPF.

## Isolation, level shifting and power states

set_isolation places clamp cells on signals leaving a domain that can be off, with a clamp value and an isolation supply that must be always on, or the cell would power down with the signal it protects. set_level_shifter requests shifters on crossings between different voltages. Power states, written with a power state table or add_power_state, list the legal combinations of supply states. From them a tool computes which crossings actually need protection: isolation where the driver can be off while the receiver is on, and shifting where the voltages differ.

## Retention and sequencing

set_retention marks registers whose state must survive power-down. A retention flop keeps its bit in a small shadow latch on an always-on supply, at a cost in area and in sleep-mode leakage, so teams often retain only state that is expensive to rebuild. The control sequence matters: isolate the outputs, save state, switch off; on wake, switch on, restore state, then release isolation. Getting the order wrong corrupts state or drives undefined values into live logic. Power-aware simulation applies the UPF so that off domains corrupt their outputs; simulation without it cannot see these bugs.

## Going deeper

IEEE Std 1801, "IEEE Standard for Design and Verification of Low-Power, Energy-Aware Electronic Systems"; Keating et al., "Low Power Methodology Manual" (2007).
