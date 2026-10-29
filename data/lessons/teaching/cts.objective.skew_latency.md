# Clock tree objectives

Clock tree synthesis (CTS) turns one clock net with thousands of sinks into a buffered network that delivers every edge on time. It runs after placement, when the flop locations are known, and before signal routing. A CTS engine is judged on a small set of numbers: latency, skew, slew and power. A tool designer needs these definitions exactly, because every algorithm in this domain optimises one of them while bounding the others, and static timing then consumes the result.

## Latency and insertion delay

The **latency** (insertion delay) of a sink is the time from the clock source to that sink's clock pin. It has two parts: the **source latency**, from the clock's origin (a PLL or a port) to the root of the tree, and the **network latency**, the sum of every buffer and wire delay along the tree path. For example, a source latency of 90 ps followed by buffers of 39, 62, 58 and 67 ps and a last wire of 10 ps gives a latency of $90 + 39 + 62 + 58 + 67 + 10 = 326$ ps. Latency itself does not appear in a setup check, because launch and capture both carry it, but large latency is costly: every stage burns power and adds on-chip variation that is not shared between branches.

## Global skew

**Skew** between two sinks is the difference of their latencies. The **global skew** of a tree is the largest latency minus the smallest over all sinks: for latencies 407, 405, 418, 386, 381 and 457 ps it is $457 - 381 = 76$ ps. Global skew is the number most CTS engines report and target, because it is cheap to compute and bounds every pairwise difference. Positive skew from launch to capture gives setup more time and takes the same amount from hold; this is the link to static timing, where skew enters each check directly.

## Local skew and why it matters more

Only flops joined by a timing path care about each other's clock. **Local skew** is the largest latency difference over sequentially adjacent pairs. If flops D, B and A are related by paths D to A, B to A and A to D, with latencies $A = 436$, $B = 445$, $D = 425$ ps, the local skew is $\max(11, 9, 11) = 11$ ps, while the global skew over a larger set might be 34 ps. A tree can therefore have large global skew and still be excellent, if unrelated sinks are the ones that differ. Useful-skew optimisation goes further and makes some local skews deliberately non-zero. A CTS engine that only minimises global skew may waste buffers equalising sinks that never talk to each other.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), the chapter on clock tree synthesis; Friedman, "Clock distribution networks in synchronous digital integrated circuits", Proceedings of the IEEE (2001).
