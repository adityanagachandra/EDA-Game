# Loop unrolling and its trade-offs

Unrolling replicates a loop body U times and runs the loop N/U times, exposing operations from several iterations in one larger basic block. It creates spatial parallelism, at the price of area and memory bandwidth, and it interacts with every other knob: resource limits, memory ports, recurrences and pipelining. Knowing where unrolling stops paying is the core of HLS design-space intuition.

## Main loop and remainder

Unrolling by U turns $N$ iterations into $\lfloor N / U \rfloor$ trips of a U-wide body. If the compiler cannot prove $N$ is a multiple of U, it must also handle the leftover $N \bmod U$ iterations, either with a separate remainder loop or by guarding each copy of the body with a condition. A loop of 105 iterations unrolled by 4 runs 26 main trips covering 104 iterations and one remainder iteration. Guards cost little but prevent some optimizations; a remainder loop costs extra control states. Full unrolling, U = N, removes the loop entirely and is only possible with a constant trip count. Partial unrolling exposes U iterations' operations to the scheduler at once, so independent iterations can execute in parallel, but only if units and memory ports exist to serve them.

## Cycles under resource limits

Model each unrolled body as taking the larger of its dependence depth $D$ and its resource bounds. With $a$ multiplications and $b$ additions per iteration, $k_m$ multipliers and $k_a$ adders, and independent iterations, one body of U copies takes $\max(D, \lceil U a / k_m \rceil, \lceil U b / k_a \rceil)$ cycles, and the loop takes that times $N / U$. With $a = 2$, $b = 3$, $D = 5$, three multipliers, three adders and U = 4, the body takes $\max(5, 3, 4) = 5$ cycles, and 32 iterations take $8 \times 5 = 40$ cycles, against 160 without unrolling. Unrolling helps exactly while the depth term dominates, because the extra copies then run in the shadow of the dependence chain. This model ignores pipelining across trips, which is usually combined with unrolling in practice.

## Where unrolling stops paying

The body stays at $D$ cycles while $\lceil U a / k_m \rceil \le D$ and $\lceil U b / k_a \rceil \le D$, that is for $U \le \min(\lfloor D k_m / a \rfloor, \lfloor D k_a / b \rfloor)$. With $D = 5$, one multiplier and $a = 2$, the multiplier allows $U \le 2$; beyond that, each extra copy adds cycles in proportion and the loop time stops falling, while area and multiplexers keep growing. The same reasoning applies to memory: U reads per trip from an array with $p$ ports cost $\lceil U / p \rceil$ cycles, so unrolling without array partitioning is often wasted. Recurrences limit unrolling too: unrolling a loop with an accumulator only lengthens the dependence chain through the accumulator unless it is reassociated. The right factor is the largest one the binding constraint can absorb, not the largest the area budget allows.

## Going deeper

Kastner, Matai and Neuendorffer, "Parallel Programming for FPGAs" (2018), chapters on loop optimizations; Coussy and Morawiec (eds.), "High-Level Synthesis: From Algorithm to Digital Circuit" (2008).
