# Loop pipelining: latency and throughput

Loop pipelining overlaps successive iterations of a loop so that a new one starts every II cycles. It is the single most effective HLS transformation for loops, because it raises throughput by up to the iteration latency divided by II, often at modest area cost. Its cost model, total cycles, iterations in flight and speedup, is simple and should be computed by reflex.

## Total cycles of a pipelined loop

With $N$ iterations, an initiation interval II and an iteration latency $L$, the last iteration starts $(N - 1)\,\mathrm{II}$ cycles after the first and takes $L$ cycles, so the loop takes $(N - 1)\,\mathrm{II} + L$ cycles. A loop of 423 iterations with II = 1 and $L = 22$ takes $422 + 22 = 444$ cycles, against $423 \times 22 = 9306$ unpipelined. For long loops the $(N - 1)\,\mathrm{II}$ term dominates, so II, not latency, decides performance. The formula assumes no stalls; variable-latency memory or stream back-pressure stretches it, and HLS reports then give a minimum and a maximum. A loop nested inside an outer loop pays the fill and drain of $L - \mathrm{II}$ cycles on every outer iteration unless the loops are flattened or the outer loop is pipelined too.

## Iterations in flight and pipeline depth

Each iteration lives for $L$ cycles and a new one starts every II cycles, so at most $\lceil L / \mathrm{II} \rceil$ iterations are in flight at once. This is the pipeline's stage count, the number of stages of II cycles that one iteration's schedule is cut into. With II = 1 and $L = 5$, five iterations overlap; with II = 4 and $L = 12$, three do. Every in-flight iteration needs its own copy of each value that is live across stage boundaries, so register count grows with depth: a value live for 6 cycles in a pipeline at II = 1 needs 6 registers, one per overlapping iteration. That is the area price of pipelining, together with the extra units needed to sustain the II, and it is why a deeply pipelined loop at II = 1 can be much larger than its unpipelined version.

## Speedup and its limit

The speedup over running iterations back to back is $\frac{N L}{(N - 1)\,\mathrm{II} + L}$. For 182 iterations with $L = 12$ and II = 4 it is $2184 / 736 = 2.97$. As $N$ grows, it approaches $L / \mathrm{II}$, here 3; for short loops the fill and drain eat into it, and a loop of 2 iterations gains almost nothing. This limit tells a designer what pipelining can buy before running the tool, and how to get more: reduce II, by adding resources, partitioning memories or breaking recurrences, or reduce the work per iteration. It also shows why a modest II on a long loop often beats aggressive unrolling: unrolling multiplies area by the unroll factor, while pipelining at II = 1 reuses one copy of the body.

## Going deeper

Kastner, Matai and Neuendorffer, "Parallel Programming for FPGAs" (2018), chapter on pipelining; the pipelining chapter of the AMD Vitis HLS user guide (UG1399).
