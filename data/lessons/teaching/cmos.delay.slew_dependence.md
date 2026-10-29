# Input slew and the limits of the RC model

The RC model assumes a step input, but every real input ramps, and a gate's delay depends strongly on how fast its input moves. This single gap is why timing libraries are tables indexed by input slew and why every timer propagates slew alongside arrival time. A tool designer must know how slew enters delay and what that means for propagation.

## Why input slew changes delay

When the input ramps instead of stepping, the pull-down transistor turns on gradually: its current grows as the input climbs past the threshold, while the pull-up is still partly on and fights it. The output therefore starts moving later relative to the input's 50% point and moves more slowly at first. The result is that delay rises with input transition time, roughly linearly over a useful range. A first-order correction adds a fitted fraction of the input slew to the step delay. Measured delay therefore depends on two quantities: input transition time and output load. That is exactly why a Liberty NLDM table has two axes, and why its rows are labelled by input transition.

## Slew propagates along a path

A gate's output transition depends on its own load and drive, and partly on its input slew. That output transition is the next gate's input slew. A timer must therefore compute, for every arc, both a delay and an output slew, and carry the slew forward with the arrival time. Resizing a buffer two stages upstream changes the slew it delivers, which changes the delay of the next gate and its output slew, and so on downstream, so a local change has non-local timing effects. Incremental timers track these dependencies when they update after an ECO. Liberty supports this with paired tables: cell_rise and rise_transition, cell_fall and fall_transition, each indexed by the same input slew and load.

## Edge cases and slew limits

Slew has surprising consequences. With a very slow input and a light load, the output can cross 50% before the input does, giving a negative delay; timers and libraries must accept this. Slow edges increase short-circuit power, because both networks conduct longer, and make the receiving gate more sensitive to noise. Delays outside the characterized slew range are extrapolated and unreliable. Libraries therefore specify a max_transition on pins, and synthesis and placement tools treat transition limits as design rules to be fixed by buffering or upsizing before timing is trusted. When a step-input estimate disagrees badly with SPICE, the usual cause is a weak driver into a small stage, where slew dominates the delay.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed.), section 4.4.6; the Liberty Reference Manual (Liberty User Guides and Reference Manual Suite), timing group and lookup tables.
