# Minimum clock period of a register path

The setup check can be turned around: instead of asking whether a path meets a given period, ask for the smallest period it meets. That number, and its inverse, the maximum frequency, is how architects and implementation teams talk about a design, and it shows directly which path limits the clock.

## The smallest period one path meets

Set the setup slack to zero and solve for $T$. With skew defined as capture latency minus launch latency, setup holds when
$$T \ge t_{cq} + d_{max} + t_{su} - \text{skew}.$$
For $t_{cq} = 50$ ps, $d_{max} = 550$ ps, $t_{su} = 60$ ps and a skew of -15 ps, the path needs $50 + 550 + 60 + 15 = 675$ ps. Negative skew, a capture clock that arrives early, raises the required period; positive skew lowers it, which is the basis of useful skew. Clock uncertainty adds directly to the bound. Each term names a lever: a faster flop, faster logic, a smaller setup time, or a later capture clock.

## The design's minimum period

A design meets a period only if every register path does, so its minimum period is the largest per-path minimum: $T_{min} = \max_p (t_{cq} + d_p + t_{su} - \text{skew}_p)$. With two paths needing 887 ps and 676 ps, the design needs 887 ps. The limiting path is the critical path, and its slack at any period $T$ equals $T - T_{min}$. Two paths with similar logic delay can rank differently once skew is counted: a path whose capture clock is late can tolerate more logic. Pipelining raises frequency by cutting the worst path's logic delay, at the cost of extra flop overhead $t_{cq} + t_{su}$ in each new stage, which limits how far it pays.

## Converting to frequency

Frequency is the inverse of the period. With the period in picoseconds, $f_{max} = 10^6 / T_{min}$ MHz, because one MHz has a period of $10^6$ ps. A minimum period of 981 ps gives about 1019.4 MHz; 1250 ps gives 800 MHz. Frequency is not linear in delay: shaving 50 ps off a 1000 ps path gains about 52.6 MHz, while the same 50 ps off a 500 ps path gains about 222 MHz. When teams quote a target frequency, the timer turns it into a period with create_clock, and all checks work in time. Reporting the margin as slack in ps rather than as MHz keeps different paths comparable.

## Going deeper

Harris and Harris, "Digital Design and Computer Architecture", the section on timing of sequential logic; Hennessy and Patterson on pipelining overhead.
