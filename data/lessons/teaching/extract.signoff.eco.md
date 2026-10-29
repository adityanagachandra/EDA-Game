# Engineering change orders

An engineering change order (ECO) is a targeted change to a nearly finished design: a logic fix found late in verification, or a timing or noise fix found at signoff. Instead of re-running synthesis and place-and-route, which would perturb everything that already closed, an ECO touches only what must change. ECOs are where extraction, timing and physical tools meet under the most time pressure, and they come in kinds with very different costs.

## Functional and timing ECOs

A functional ECO changes what the logic computes: it fixes a bug, such as an inverted enable or a wrong reset value. Equivalence checking against the corrected RTL must then pass, while equivalence against the old netlist fails in exactly the intended way. A timing ECO changes cells or wires without changing function: upsizing or downsizing a cell, swapping its threshold-voltage flavour, inserting a buffer or delay cell, or rerouting a net. It must remain logically equivalent to the pre-ECO netlist. Both are preferred to a re-run of the flow because weeks of timing, SI and physical closure elsewhere are preserved. ECO tools work at the netlist and placement level, then incrementally route, extract and re-time.

## Metal-only ECOs and spare cells

After tape-out, masks are expensive, and wafers are often held after the base (transistor) layers are processed. A metal-only ECO changes only metal and via layers, so the base masks are reused and held wafers can be finished with the new metal masks. It cannot add transistors, so designers scatter spare cells (NANDs, NORs, inverters, flip-flops, with inputs tied off) across the die before tape-out. A fix is mapped onto nearby spares and wired in. Spare flops often have their clock pins prewired so the clock tree need not change. Some flows use gate-array-style filler cells that can be configured as logic by metal alone, making more resources available.

## Side effects and re-verification

An ECO is local in intent but not in effect. A new buffer adds load to its driver and coupling to its neighbours; a rerouted wire changes the capacitance of every wire beside it; an upsized cell speeds short paths through it and can break hold. So after an ECO, extraction must cover the touched nets and their geometric neighbours, timing and SI must be re-run, and EM, IR, DRC, LVS and equivalence must be re-checked. A functional ECO is proven with formal equivalence against the corrected RTL, not by a simulation of the failing test, since only equivalence covers every input. The best ECO is the smallest change that passes all of these, which is why ECO tools search for minimal patches.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter on timing closure; the OpenROAD documentation for its ECO and repair commands.
