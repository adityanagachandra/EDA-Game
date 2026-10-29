# Cell characterization

Characterization turns transistor netlists and foundry device models into Liberty: it simulates every arc of every cell at every grid point and corner, measures delays, transitions, constraints, power and capacitances, and writes the tables. It is the bridge between SPICE and every digital tool, and its definitions and cost determine how much a library can be trusted and how fast it can be rebuilt.

## The simulation count

For each cell, the characterizer enumerates arcs from the function and pins, chooses side-input states that sensitize each arc, and runs a transient simulation per grid point and output edge. The count multiplies quickly:
$$N = \text{cells} \times \text{arcs per cell} \times 2 \times n_1 n_2 \times \text{corners}.$$
A library of 400 cells averaging 4 arcs, with 7 by 7 tables, both edges and 9 corners, needs $400 \times 4 \times 2 \times 49 \times 9 = 1{,}411{,}200$ transients for delay alone. Constraint arcs add many more, because each setup or hold value comes from a bisection search with several simulations per point. This is why characterization runs on large compute farms, why libraries are characterized at a limited set of corners, and why faster methods are an active area.

## Measuring delay and transition

Each transient applies an input ramp and records the output. Delay is measured between threshold crossings defined in the library header, typically 50% of the swing on input and output (input_threshold_pct and output_threshold_pct). Transition is measured between lower and upper thresholds, such as 20% and 80% or 10% and 90%. Example: an input ramping from 0 to $V_{DD}$ between 10 and 70 ps crosses 50% at 40 ps; an output falling linearly between 50 and 130 ps crosses 50% at 90 ps, so the delay is 50 ps. That output spends 60% of its 80 ps swing between the 20% and 80% points, so its 20-80 transition is 48 ps. The load is a capacitor at the table's load value, and the input slew is the row's index value.

## Thresholds, derating and conventions

Libraries differ in how they store transitions. Some store the measured threshold-to-threshold time; others store the equivalent full-swing ramp time, the measured time divided by the fraction of the swing between thresholds, with the convention recorded through attributes such as slew_derate_from_library and the threshold percentages. A 30-to-70 measurement of 20 ps corresponds to a full-swing ramp of $20/0.4 = 50$ ps. A tool that applies one convention to a library written in another is off by that factor on every slew, which then propagates through every lookup. Characterization also measures pin capacitances, leakage per state, internal energy and constraints, and its own correctness is verified by comparing selected points against independent SPICE runs.

## Going deeper

The Liberty Reference Manual, library-level threshold and slew attributes; Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 3.
