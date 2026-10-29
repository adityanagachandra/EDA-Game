# Operation latency and the clock period

A scheduler counts in control steps, but operations have delays in nanoseconds. Converting one into the other, given the target clock period and the characterized delay of each operator, is the first quantitative decision of HLS. It decides whether a multiplication is a one-cycle or a three-cycle operation, which in turn sets the critical path, the unit occupancy and the register count of everything downstream.

## From nanoseconds to cycles

Each operator type has a characterized combinational delay for its width on the target technology: a 32-bit adder might take 1.2 ns, a 32-bit multiplier 6 ns. If the result must be registered before use and register overhead is ignored, an operation with delay $d$ at clock period $T$ needs $\lceil d / T \rceil$ cycles. A 7.3 ns multiplier at a 3 ns clock needs $\lceil 2.43 \rceil = 3$ cycles. The ceiling is the cost of quantization: the third cycle is mostly idle, but it still costs a whole cycle. In practice the usable part of the period is $T$ minus the register's clock-to-q and setup time, plus a margin for clock uncertainty, so the tool divides by that usable time instead. The operator library that supplies these delays is one of the most important inputs an HLS engine has.

## Latency of a registered chain

When every result is registered and nothing is chained, a chain of dependent operations takes the sum of each operation's cycle count. Delays of 4, 2.8 and 4.9 ns at a 4 ns clock cost 1, 1 and 2 cycles, so the chain takes 4 cycles, or 16 ns, although its total combinational delay is only 11.7 ns. The wasted 4.3 ns is the quantization loss of three separate roundings. This is why schedulers chain short operations into one cycle when they fit, and why HLS reports often show a latency in cycles that looks larger than the raw arithmetic would suggest. A multicycle operation also occupies its unit for several cycles, unless the unit is pipelined, which changes how many operations a unit can serve per cycle.

## Choosing the clock period

A faster clock is not automatically better. The total time of a chain is its cycle count times the period, and the rounding interacts with the period. With delays 9, 7.5 and 5.6 ns, a 2.5 ns clock gives $4 + 3 + 3 = 10$ cycles, 25 ns, while a 5 ns clock gives $2 + 2 + 2 = 6$ cycles, 30 ns. Another chain might favour the slower clock. Beyond rounding, a faster clock adds pipeline registers, area and power, and raises the register overhead's share of each cycle. Designers usually fix the period from the surrounding system and let HLS fit the operations to it; a design-space explorer may sweep the period as one more knob. Either way, the latency and throughput numbers an HLS report shows depend on the period at which the tool was asked to schedule.

## Going deeper

De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 5; the operator characterization sections of the AMD Vitis HLS user guide (UG1399).
