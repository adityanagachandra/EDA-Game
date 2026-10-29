# Parallel simulation: work and span

A large SoC simulates slowly on one core, so simulators try to use many. A partitioned design is a task graph: each partition is work, and data dependences between partitions are edges. Two numbers bound what parallelism can deliver, the total **work** and the critical-path **span**. Their ratio and Brent's bound tell a simulator designer how much speedup is available before paying for synchronization, and why fine-grained parallel logic simulation has a long history of disappointing results.

## Work and span of a cycle's task graph

Split one simulated cycle into tasks (partitions of the logic) with known costs, and draw an edge from task $u$ to task $v$ when $v$ reads a value $u$ produces in the same cycle. The **work** $W$ is the sum of all task costs: the time one thread needs. The **span** $S$ is the cost of the most expensive path through the graph: the time an unlimited number of threads needs, since tasks on a path must run in order. Example: T1 costs 4 us, T2 costs 3 us, T3 (after T1) costs 5 us, T4 (after T2 and T3) costs 2 us. Then $W = 14$ us. The paths are T1, T3, T4 with cost 11 us and T2, T4 with cost 5 us, so $S = 11$ us.

## Speedup bounds

With $P$ threads, the time per cycle $T_P$ satisfies $T_P \ge \max(W/P, S)$: you cannot do the work faster than $P$ at a time, and you cannot shorten the critical path. The best possible speedup is therefore $\min(P, W/S)$. The ratio $W/S$ is the **parallelism**: in the example, $14 / 11 \approx 1.27$, so even a thousand threads give at most 1.27 times speedup. In the other direction, Brent's theorem says a greedy scheduler, which never idles a thread while a task is ready, achieves $T_P \le W/P + S$. For the example with $P = 4$: $T_4 \le 14/4 + 11 = 14.5$ us, a weak bound here because the span dominates. Good partitioning aims for many independent tasks of balanced cost, so that $W/S$ is much larger than $P$.

## Why synchronization dominates

The bounds ignore communication. Every simulated cycle ends with a barrier: all threads must finish before registers update, and values crossing partitions must be exchanged. If a barrier costs a few microseconds and a cycle's work per thread is comparable, synchronization eats the gain. Fine-grained schemes, one task per gate or per event, make it worse: a gate evaluation costs nanoseconds, far less than any synchronization. That is why practical parallel simulators, including Verilator's multithreaded mode, merge logic into coarse tasks, balance them by estimated cost, and minimize the values that cross between threads, which is a graph-partitioning problem with a balance constraint. The scaling limit is then the design's own structure: a long serial dependency chain in one cycle caps speedup no matter how many cores are added.

## Going deeper

Blumofe and Leiserson, "Scheduling multithreaded computations by work stealing", Journal of the ACM (1999); Brent, "The parallel evaluation of general arithmetic expressions", Journal of the ACM (1974).
