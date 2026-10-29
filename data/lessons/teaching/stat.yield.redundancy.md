# Redundancy and repairable yield

Large memories would rarely work if every cell had to be perfect. Instead they carry spare rows and columns, and a die whose faults can be mapped onto spares is repaired at test. Redundancy turns yield from the probability of zero faults into the probability of at most as many faults as spares, and sizing it is a direct computation with the Poisson or binomial distribution. The same arithmetic applies to spare cores and spare lanes.

## Yield with spares

Suppose the number of faults in a memory array is Poisson with mean $\lambda$, each fault is repairable by one spare, and the spares themselves are fault-free. With $R$ spares the array is good when $K \le R$:
$$Y_R = \sum_{k=0}^{R} \frac{e^{-\lambda}\lambda^k}{k!}.$$
With $\lambda = 1.5$, no spares gives $e^{-1.5} = 22.3$ percent; two spares give $e^{-1.5}(1 + 1.5 + 1.125) = 80.9$ percent; four give 98.1 percent. A handful of spares rescues an array whose unrepaired yield is poor. The assumptions are optimistic: in reality spares can be faulty, faults cluster, and a single defect can hit several cells that need different spares.

## Sizing the redundancy

The design question runs backwards: how many spares reach a target yield? Accumulate the Poisson cumulative distribution until it reaches the target. With $\lambda = 1.5$ and a target of 99 percent, $P(K \le 4) = 0.981$ and $P(K \le 5) = 0.9955$, so five spares are needed. Each spare costs area, which slightly raises $\lambda$, and repair logic and fuses cost more area and test time, so the optimum is a trade between unrepaired loss and overhead. Repair allocation itself, choosing which spare rows and columns cover a fault map, is an NP-complete covering problem solved with branch-and-bound or heuristics in built-in repair analysis.

## Block-level redundancy

At a coarser grain, a chip needs $n$ working blocks (cores, sub-arrays, lanes) and is built with $n + s$. If each block is defective independently with probability $p$, the chip is good when at most $s$ blocks fail:
$$Y = \sum_{k=0}^{s} \binom{n+s}{k} p^k (1-p)^{n+s-k}.$$
For $n = 8$, $s = 1$ and $p = 0.05$: without a spare, $0.95^8 = 66.3$ percent; with one spare among nine, $0.95^9 + 9(0.05)(0.95^8) = 0.630 + 0.299 = 92.9$ percent. Selling parts with some blocks disabled applies the same arithmetic. The model again assumes independent block failures, which clustering violates, usually in the direction of making spares slightly less valuable than predicted.

## Going deeper

Koren and Koren, "Defect tolerance in VLSI circuits: techniques and yield analysis," Proceedings of the IEEE, 1998; Kuo and Fuchs, "Efficient spare allocation for reconfigurable arrays," IEEE Design and Test, 1987.
