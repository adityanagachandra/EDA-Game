# Hardware-software partitioning

HLS is usually applied to part of a program: the kernel that dominates run time is moved to an accelerator while the rest stays in software. Deciding what to move, and whether the move pays once data transfer is counted, is hardware-software partitioning. Amdahl's law and a transfer-cost term give most of the answer before any synthesis is run.

## New run time after offloading

Let the software run time be $T$, the fraction spent in the candidate kernel $f$, the accelerator's speedup on that kernel $s$, and the added time for moving data to and from the accelerator $c$. The new run time is $T' = (1 - f) T + f T / s + c$. For $T = 440$ ms, $f = 0.7$, $s = 58$ and $c = 70.8$ ms, the software part stays at 132 ms, the kernel drops from 308 ms to 5.31 ms, and the transfer adds 70.8 ms, for 208.11 ms. Transfer often dominates: here it costs more than thirteen times the accelerated kernel itself. Overlapping transfer with computation, by streaming or double buffering, removes part of $c$, and keeping data resident on the accelerator across calls removes more.

## Overall speedup and Amdahl's limit

The overall speedup is $T / T'$. With $T = 780$ ms, $f = 0.83$, $s = 19$ and $c = 58.3$ ms, $T' = 132.6 + 34.07 + 58.3 = 224.97$ ms and the speedup is 3.47. However fast the accelerator, the speedup can never exceed $1 / (1 - f)$, here 5.88, because the software part remains. This is Amdahl's law, and it is the first filter in partitioning: a kernel covering 50% of run time can never give more than 2x overall. It also guides profiling: candidates should be ranked by the fraction of time they cover and by how much data they touch per unit of compute, since kernels with high arithmetic intensity amortize transfer best. A kernel that does little work per byte moved may need a speedup it can never reach.

## The break-even speedup

An accelerator pays off only if the kernel time it saves exceeds the transfer it adds: $f T / s + c < f T$, so $s > \frac{f T}{f T - c}$. For $T = 110$ ms, $f = 0.41$ and $c = 5$ ms, the kernel takes 45.1 ms in software, and the accelerator must run it at least $45.1 / 40.1 = 1.12$ times faster just to break even. If $c \ge f T$, no accelerator can help, because transfer alone costs more than the kernel. Partitioning tools formalize this as an optimization over many candidate kernels, with costs for area, transfer and synchronization, solved by ILP or heuristics; design-space exploration then tunes each chosen kernel. The same reasoning applies on-chip, between a processor and an accelerator sharing a cache.

## Going deeper

De Micheli and Gupta, "Hardware/software co-design" (Proceedings of the IEEE, 1997); Hennessy and Patterson, "Computer Architecture: A Quantitative Approach", the discussion of Amdahl's law.
