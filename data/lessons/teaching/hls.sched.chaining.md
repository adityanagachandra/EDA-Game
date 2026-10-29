# Operation chaining

Chaining places several dependent operations in one clock cycle when their combined delay fits, instead of registering every result. It cuts latency and registers, and it is why an HLS schedule is not simply the DFG depth in cycles. The scheduler must know operator delays, register overhead and the clock period, and it must stop chaining before the cycle overflows.

## Usable time in one cycle

A register-to-register path must fit in the clock period minus the register overhead: clock-to-q of the launching register, setup of the capturing one, and any margin for clock uncertainty. With a 10 ns clock and 0.5 ns of overhead, 9.5 ns of logic fits per cycle. Operations chained in one cycle add their delays: an adder of 1.6 ns followed by a multiplier of 5.7 ns uses 7.3 ns, and a further 1.8 ns adder brings it to 9.1 ns, still within 9.5. A fourth operation of 5.4 ns would exceed it and must start a new cycle. Counting how many leading operations of a chain fit in the first cycle is the basic chaining decision. The delays come from the same operator characterization that converts multicycle operations to cycle counts.

## Packing a chain into cycles

For a single chain of operations in fixed order, greedy packing is optimal: fill each cycle with as many next operations as fit, then start a new one. The delays 0.9, 1.0, 2.8, 1.4 and 1.7 ns with 4.7 ns usable per cycle pack as $0.9 + 1.0 + 2.8 = 4.7$ in the first cycle and $1.4 + 1.7 = 3.1$ in the second, two cycles instead of five. Without chaining, each operation takes its own cycle, so the saving is the difference, three cycles here. In a general DFG the scheduler chains along the edges, tracking for each operation the arrival time of its inputs within the current cycle, which is a small static timing analysis performed during scheduling. SDC schedulers encode the same rule as difference constraints between pairs of operations whose combined delay exceeds the period.

## When chaining hurts

Chaining trades clock period headroom for latency. Every chained path is a long combinational path, so the estimate must be accurate: operator delays in HLS libraries are pre-layout numbers, and wires, fanout and multiplexers added by binding lengthen the real path. An optimistic chain that fits by 0.1 ns in the HLS model can fail timing after placement, which is why tools apply a clock uncertainty margin. Chaining also interacts with sharing: an operation chained behind a shared unit inherits that unit's input multiplexer delay. And with no registers between chained operations there is no place to retime later. When delays are uncertain, a tool may chain conservatively, for example only short operations such as additions and logic, and leave multiplications registered.

## Going deeper

De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 5; Cong and Zhang, "An efficient and versatile scheduling algorithm based on SDC formulation" (DAC, 2006), section on operation chaining.
