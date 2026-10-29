# Delta cycles

Real gates have delay, but RTL is usually written with none: continuous assignments and always blocks update in zero simulated time. A simulator still has to order those updates causally. It does so with **delta cycles**, successive passes at the same simulated time, each applying the changes the previous pass produced. Delta cycles are how zero-delay RTL settles, where zero-time glitches come from, and what a simulator reports when a combinational loop never settles.

## Ordering zero-delay updates in delta steps

When a net changes at time $t$ and the logic reading it has zero delay, the resulting change is also due at time $t$, but it must happen after its cause. The simulator therefore splits time $t$ into delta steps $t + \delta_1, t + \delta_2, \ldots$: in each step it applies pending updates, wakes the processes and gates that read the changed nets, evaluates them, and queues their output changes for the next delta. Simulated time does not advance until a delta produces no new changes. Delta steps are invisible in a waveform viewer, which shows only the final value at each time, but they are real work: every delta is another round of evaluations.

## Counting deltas to settle

Consider the chain `assign w1 = ~a; assign w2 = ~w1; assign w3 = ~w2;`. If the testbench changes $a$ at 10 ns, $a$ changes in delta 1, $w_1$ in delta 2, $w_2$ in delta 3 and $w_3$ in delta 4. Four deltas contain a change, then time advances. In a general zero-delay netlist, the number of deltas is the length of the longest chain of gates whose outputs actually change, plus one for the input. It is not the depth of the whole netlist: if an AND gate with a 0 on its other input absorbs the change, the propagation stops there. To count, trace the change level by level, as for unit delay, with each unit of delay read as one delta.

## Zero-time glitches and loops

When reconvergent paths have different lengths, their changes arrive at a gate in different deltas. Take $y = a \oplus w$ with $w = \lnot\lnot a$ through two inverters. When $a$ rises in delta 1, $y$ is re-evaluated with the new $a$ and the old $w$ and rises in delta 2; $w$ changes in delta 3, and $y$ falls back to 0 in delta 4, all at one simulated time. The glitch is invisible in time but triggers any always block sensitive to $y$, which is a classic source of spurious activity and of races when testbench code samples during deltas. A combinational loop with odd inversion never settles: deltas continue forever at one time. Simulators guard against this with an iteration limit and report a zero-delay oscillation.

## Going deeper

IEEE Std 1800-2017, clause 4 (scheduling semantics); the VHDL LRM, IEEE Std 1076, on the simulation cycle and delta delay.
