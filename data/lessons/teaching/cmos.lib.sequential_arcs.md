# Sequential constraint arcs in Liberty

A flip-flop's Liberty description has delay arcs from the clock to the outputs and constraint arcs from the clock to the data and control inputs. Setup and hold are stored as tables indexed by the data and clock transitions, and a timer turns them into required times. Reading these tables correctly is the link between transistor-level characterization and every setup and hold slack in a report.

## Constraint tables

A setup or hold check is a timing group on the data pin D with related_pin CK and timing_type setup_rising or hold_rising for a rising-edge flop. Its tables, rise_constraint and fall_constraint for a rising or falling D, are indexed by the transitions of the constrained pin (D) and the related pin (CK). Which variable is index_1 is set by the table template, so a reader must check the template rather than assume. Example: index_1 = data transition 10, 30, 90 ps and index_2 = clock transition 10, 40, 160 ps, with the row for 30 ps reading 24, 20, 8. At a 30 ps data transition and a 40 ps clock transition, the setup time is 20 ps. Setup typically grows with data slew and shrinks with clock slew.

## Setup slack from a constraint

The timer uses the looked-up setup time to set the required time at D. For a capture edge one period after launch, with clock latency $L$ at CK, setup $t_{su}$ and data arrival $A$:
$$\text{slack}_{setup} = (T + L - t_{su}) - A.$$
With a 500 ps period, 60 ps of clock latency, a 20 ps setup time and data arriving at 515 ps, the required time is 540 ps and the slack is 25 ps. The clock transition used for the lookup is the slew at the CK pin, so a slow clock tree changes the setup time itself, not just the clock arrival. Clock-to-q delay, from the rising_edge arc on Q, starts the next path.

## Hold slack and other checks

Hold compares the earliest data change with the same clock edge:
$$\text{slack}_{hold} = A_{min} - (L + t_h).$$
With 40 ps of clock latency, a 6 ps hold time and the earliest data change at 58 ps, hold slack is 12 ps. Hold tables can contain negative values, which mean data may change slightly before the edge; a reader must preserve the sign. Beyond setup and hold, flops with asynchronous set or reset have recovery and removal arcs (timing_type recovery_rising, removal_rising), which play the setup and hold roles for the release of the reset, and min_pulse_width constraints on clock and reset pins. All of them come from the characterization searches that found the cell's failure boundaries.

## Going deeper

The Liberty Reference Manual, timing_type constraint values; Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapters 7 and 8.
