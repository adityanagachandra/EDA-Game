# Interconnect and multiplexer cost

Sharing a functional unit or a register is not free: every distinct source that feeds a shared input port needs a multiplexer input. In FPGA designs these multiplexers can cost as much area as the units they enable sharing for, and they add delay to every path through the shared unit. Interconnect-aware binding is where HLS quality is often won or lost.

## Sources at a port

After binding, each input port of a unit receives operands from some set of registers or units across the steps in which it is used. A port fed from one source needs only a wire; a port fed from $k$ distinct sources needs a $k$-input multiplexer, whose select signal comes from the controller. Count the distinct sources, not the operations: if adder A1 executes $R_2 + R_5$, $R_1 + R_5$, $R_1 + R_3$ and $R_4 + R_1$, its left port reads $R_2$, $R_1$ and $R_4$, three distinct registers, so it needs a 3-input multiplexer, and its right port reads $R_5$, $R_3$ and $R_1$, also three. Register inputs are counted the same way: a register written by several units needs a multiplexer on its data input. The total multiplexer area grows with every additional distinct connection.

## Counting two-input multiplexer equivalents

A $k$-input multiplexer can be built from $k - 1$ two-input multiplexers, so a common cost measure is $\sum_{\text{ports}} (k_p - 1)$ over ports with $k_p \ge 1$ sources. For an adder executing $R_4 + R_6$, $R_3 + R_5$ and $R_2 + R_5$, the left port has 3 sources and the right port 2, so the cost is $(3 - 1) + (2 - 1) = 3$ two-input multiplexers. On FPGAs the real cost is in LUTs and depends on width: a 32-bit 4-input multiplexer costs 32 times a 1-bit one. Delay matters too: each level of multiplexing lengthens every path through the shared unit, which can cancel the clock-period gain the scheduler assumed. This is why HLS tools sometimes refuse to share small units such as adders, whose cost is close to that of the multiplexers that sharing them would need.

## Using commutativity

For commutative operations, the binder may swap operands, which changes which port each source feeds. The adder executing $R_3 + R_1$, $R_5 + R_2$, $R_1 + R_6$ and $R_3 + R_2$ costs $(3 - 1) + (3 - 1) = 4$ as written, because the left port reads $R_3, R_5, R_1$ and the right $R_1, R_2, R_6$. Swapping the third addition to $R_6 + R_1$ makes the left port read $R_3, R_5, R_6$ and the right $R_1, R_2$, for a cost of $2 + 1 = 3$. With $n$ commutative operations there are $2^n$ orientations, so binders use heuristics or fold the choice into a weighted bipartite matching between sources and ports. The same idea extends to choosing which register holds which value: binding values that feed the same port into the same register removes a multiplexer input entirely.

## Going deeper

De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 6, on interconnect and binding; Gajski, Dutt, Wu and Lin, "High-Level Synthesis" (1992), chapter on data-path synthesis.
