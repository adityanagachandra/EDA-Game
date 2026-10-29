# Bank conflicts and the II they force

Partitioning only helps if the accesses of one iteration spread over the banks. When two accesses in the same cycle hit the same bank, they conflict, and the schedule must spread them over extra cycles. Predicting bank conflicts from the access pattern is the computation an HLS memory partitioner performs, and it explains why a partitioned design can still miss its II.

## Which bank an access hits

With cyclic partitioning into $B$ banks, the access $a[i + o]$ touches bank $(i + o) \bmod B$. In the iteration with $i = 1$ and $B = 3$, the reads $a[i+1], a[i+2], a[i+3], a[i+4]$ and $a[i+8]$ touch elements 2, 3, 4, 5 and 9, which are in banks 2, 0, 1, 2 and 0. Two reads hit bank 2 and two hit bank 0. The pattern matters, not the absolute index: changing $i$ by one adds one to every bank number modulo $B$, which renames the banks but keeps the same grouping. So the conflict structure of a loop with constant offsets is determined by the offsets modulo $B$, and a partitioner can analyze it once, symbolically, for all iterations. For block partitioning the bank is the index divided by the bank size, and the pattern does change with $i$, which makes analysis harder.

## The busiest bank

Count the accesses per bank in one iteration and take the maximum. For offsets 0, 5, 6 and 7 with $B = 4$, the residues are 0, 1, 2 and 3, so every bank receives one access: a conflict-free pattern. For offsets 3, 4, 5 and 7 with $B = 2$, the residues are 1, 0, 1 and 1, so one bank receives three. The maximum per bank is the quantity that matters, because the busiest bank determines how many cycles the iteration's accesses need. A partitioner chooses $B$ to make the residues of the offsets as evenly spread as possible; for $k$ distinct consecutive offsets, $B \ge k$ guarantees one access per bank, but irregular offsets can conflict even with many banks, for example offsets 0 and 4 with $B = 4$.

## From conflicts to II

If the busiest bank receives $m$ accesses per iteration and each bank has $p$ ports, those accesses need $\lceil m / p \rceil$ cycles, and every iteration repeats them, so $\mathrm{II} \ge \lceil m / p \rceil$. Three accesses to one single-port bank force II = 3; with dual-port banks, II = 2. This is a ResMII bound in which each bank is its own resource. The fixes are a larger or different factor, a different partitioning dimension, or rewriting the code so that values are read once and passed through registers, as a shift register does for sliding windows. HLS reports phrase the result as an II violation on a memory port; the designer's job is to trace it back to the offsets and the bank function.

## Going deeper

Cong, Jiang, Liu and Zou, "Automatic memory partitioning and scheduling for throughput and power optimization" (ICCAD, 2009); Wang, Li and Cong, "Theory and algorithm for generalized memory partitioning in high-level synthesis" (FPGA, 2014).
