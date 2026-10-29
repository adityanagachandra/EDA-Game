# Electrical rule checking

Electrical rule checks look for circuits that are connected as drawn but electrically wrong: floating gates, untapped wells, two drivers on one net, devices exposed to voltages beyond their rating. They run on the extracted netlist, often together with LVS, and they catch mistakes that a schematic comparison cannot, because the schematic may contain the same mistake.

## Floating gates

A gate net is floating if it reaches only gate terminals and is not a port of the cell: nothing drives it. In a netlist of an inverter, a NAND2 and a NOR2 whose inputs are A to E, any input net not declared as a port is floating. A floating gate settles at an undefined voltage and can leave both transistors of a stage partly on, burning static current. The check is a pass over nets: collect each net's terminal kinds, then flag nets with only gate terminals and no port or driver.

## Well taps

Every n-well must be tied to its supply through an n-plus tap, and the substrate tied to ground through p-plus taps. An untapped well floats, which shifts the PMOS thresholds through the body effect and raises the risk of latch-up, the parasitic thyristor formed by the n-well and substrate. The ERC check tests, for each well region, whether at least one tap lies inside it; a stricter rule bounds the distance from any transistor to the nearest tap. Standard-cell rows usually rely on tap cells placed at a fixed pitch, so a missing tap cell shows up here.

## Other electrical rules

Further ERC checks include multiple drivers on one net (contention), gates tied directly to a supply rather than through tie-high or tie-low cells (oxide stress), signals crossing between power domains without a level shifter, thin-oxide devices connected to high-voltage nets, and dangling nets that connect to a single terminal. Voltage-aware checks propagate the possible voltage of each net from the supplies and compare it with each device's rating. ERC findings are often warnings rather than errors, and decks distinguish the two so that signoff is not blocked by benign cases.
