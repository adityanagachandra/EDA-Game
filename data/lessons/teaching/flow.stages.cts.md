# Clock tree synthesis

Every flip-flop needs the clock, and it needs it at nearly the same time as its neighbours. Clock tree synthesis (CTS) builds the buffered network that distributes the clock from its source to thousands or millions of sinks, controlling skew, slew and power. It runs after placement, when sink locations are known, and before routing, which must honour the clock wires it creates.

## Latency and skew

The insertion delay, or latency, of a sink is the time from the clock source to that sink's clock pin. Skew between two sinks is the difference of their latencies; global skew is the maximum minus the minimum over all sinks. If four sinks have latencies 412, 398, 425 and 405 ps, global skew is $425 - 398 = 27$ ps. CTS aims for small skew between flops that talk to each other, controlled slew at every sink, and low power, since the clock net switches every cycle and often consumes a large share of dynamic power.

## Building the buffered tree

A single driver cannot drive a million sinks, so CTS builds a tree of buffers in which each buffer drives a limited number of loads. If each buffer may drive at most $k$ loads, the leaf level needs $\lceil N/k \rceil$ buffers to serve $N$ sinks, the level above needs $\lceil \cdot/k \rceil$ of those, and so on up to one buffer at the root. For $N = 1000$ and $k = 8$: 125, then 16, then 2, then 1, a total of 144 buffers in 4 levels. Real CTS also balances wire lengths, for example with H-tree structures or by clustering sinks geometrically and balancing delay at each merge point.

## Skew and timing

Skew changes timing checks. For a path from launch flop A to capture flop B, if B's clock arrives later than A's by $\delta$, setup slack grows by $\delta$ and hold slack shrinks by $\delta$. Skew is therefore not purely bad: useful-skew optimization delays some clocks on purpose to borrow time for critical paths. But skew that CTS did not intend, from a bad tree or from on-chip variation, can create hold violations that appear only after CTS. That is why the flow re-runs timing analysis and hold fixing immediately after CTS.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 7; the OpenROAD TritonCTS documentation.
