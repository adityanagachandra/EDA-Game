# Insertion-delay balancing across sinks

Not every clock sink is a simple flop pin. Macros such as memories have their own internal clock trees, IO registers sit far from the core, and some sinks hang off generated clocks. CTS must balance the latency at the registers that actually exchange data, not only at the pins it sees. It does so by targeting a common latency and padding early groups with delay. Understanding this is necessary to read a CTS report and to set the engine's targets.

## Register-level latency

A macro's clock pin is not where its registers capture. A memory with 150 ps of internal insertion delay from its clock pin to its internal registers, whose pin sees 280 ps, has a register-level latency of 430 ps. If the core flops see 359 ps and the IO flops 298 ps, the register-level skew is $430 - 298 = 132$ ps; balancing pins alone would report $359 - 280 = 79$ ps and be wrong. CTS tools therefore accept an insertion-delay value for macro clock pins, often taken from the macro's timing model, and balance the sum.

## Padding early groups

The usual target is the largest register-level latency, since delay can only be added. A group at latency $l$ below a target $T^\ast$ with tolerance $\tau$ needs $\lceil (T^\ast - \tau - l)/d \rceil$ delay cells of delay $d$. With a target of 397 ps, a tolerance of 5 ps and 25 ps cells, a group at 251 ps needs $\lceil 141/25 \rceil = 6$ cells. In practice the padding goes high in the tree, where one inserted cell delays a whole subtree, rather than at each sink.

## The cost of balancing

Padding is not free. Every delay cell toggles each cycle and adds power. Every added stage is a new source of on-chip variation that is not shared with the other groups, so balanced nominal latency can still mean skew in silicon. For four groups at 369, 256, 258 and 368 ps with a 20 ps tolerance and 25 ps cells, eight cells are inserted. When the gap comes from one large macro, a better fix may be to move the macro's clock entry point, use a faster path to it, or exploit useful skew, so that the rest of the design is not slowed to the macro's latency.

## Going deeper

Kahng et al., "VLSI Physical Design" (2nd ed., 2022), the clock chapter; the OpenROAD documentation for the CTS module.
