# Level shifters, isolation and retention cells

Designs with several supply voltages or with blocks that power down need cells that ordinary logic libraries do not provide: level shifters between voltage domains, isolation cells at the outputs of switchable domains, retention flops that keep state through power-down, and always-on buffers. Power-intent tools insert them from a UPF or CPF description, and timing, power and verification tools must recognize each one. Their circuits explain the rules those tools enforce.

## Level shifters

A signal crossing from a low-voltage domain to a higher one cannot drive the receiver directly. A 0.6 V logic high applied to an inverter on a 0.9 V supply gives its pMOS a gate-source voltage of -0.3 V, which may exceed the pMOS threshold, so the pMOS stays partly on: static current flows and the output low is degraded. A low-to-high level shifter uses a cross-coupled pMOS pair on the high supply, with the input and its complement driving nMOS pull-downs, so that the high supply's pMOS devices are switched by the shifter's own nodes. It needs both supplies, and Liberty marks it with is_level_shifter and per-pin related power pins. A high-to-low crossing is easier: an ordinary buffer on the lower supply can accept a higher input, though libraries still provide dedicated cells for clarity and checking.

## Isolation and always-on cells

When a domain is powered down, its outputs float at undefined voltages. A floating input can sit near the switching threshold of the receiving gate, turning on both of its networks and producing unknown logic. An isolation cell, controlled by an isolation enable, clamps such a signal to a defined value, logic 0, logic 1 or the last value latched, while the source is off. It must be powered from a supply that stays on, typically placed in or powered by the receiving domain, or it would turn off exactly when needed. The clamp value is a functional decision: it must be the value the receiver treats as inactive. Always-on buffers serve a related need: they sit inside a switchable region but draw power from a secondary always-on supply pin, so control signals can be routed through a region that may be off.

## Retention flops and power-intent checking

A retention flip-flop adds a shadow latch, powered by the always-on supply, to an ordinary flop. Before power-down a save signal copies the state into the shadow latch; after power-up a restore signal copies it back. This lets a block resume without recomputing state, at the cost of area and always-on leakage, so only state worth keeping is retained. Because these cells change behaviour with supply state, verification goes beyond logic: power-aware simulation corrupts the outputs of off domains, and static power-intent checkers verify that every crossing has the right isolation and level shifting, that control signals come from always-on logic, and that each special cell's supply pins connect to the right nets. Liberty attributes and the power-intent file together make these checks possible.

## Going deeper

IEEE 1801 (Unified Power Format); Keating, Flynn, Aitken, Gibbons and Shi, "Low Power Methodology Manual" (2007).
