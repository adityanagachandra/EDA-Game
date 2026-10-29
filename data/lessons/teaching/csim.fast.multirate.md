# Fast SPICE partitioning and multirate

Fast SPICE engines simulate circuits with millions of transistors, far beyond what one global Newton solve per time point allows. They partition the circuit into weakly coupled blocks, step each block at its own rate, skip blocks that are not changing, and solve small matrices instead of one huge one. Each trick has a measurable payoff and a measurable risk, and estimating both is how a tool designer chooses among them.

## Multirate time stepping

A single-rate simulator steps the whole circuit at the step the fastest block needs. If blocks need steps $h_1, \ldots, h_B$ over a run of length $T$, single-rate work is $B \cdot T/h_{min}$ block solves, while multirate work is $\sum_i T/h_i$. With three blocks needing 10, 20 and 100 ps, the speedup is $3/(1 + 0.5 + 0.1) = 1.875$. Real circuits, with a few fast clock nets and many slow analog or idle digital nodes, have far larger spreads. The cost is coupling: a slow block's inputs from a fast block must be interpolated, and if the coupling is strong, interpolation errors feed back and the blocks must be iterated or re-solved.

## Latency and event-driven evaluation

Most of a large digital or memory circuit is idle at any instant. A latent block, whose inputs have not changed beyond a tolerance, need not be re-evaluated: its previous solution still holds. If $N$ partitions exist, $a$ are active per step, and checking a latent one costs a fraction $c$ of a full evaluation, the speedup is $N/(a + (N - a)c)$. With 2000 partitions, 20 active and $c = 0.05$, that is $2000/(20 + 99) = 16.8$. The check itself becomes the bottleneck as activity falls, which is why event-driven engines track changes by events rather than scanning every partition at every step.

## Partitioned matrix solves

Partitioning also shrinks the linear algebra. If each matrix were factored densely, with cost proportional to size cubed, a 3000-unknown circuit split into 12 blocks of 250 plus a 50-unknown interface system would cost $12 \cdot 250^3 + 50^3$ instead of $3000^3$, a ratio of about 144. Real circuit matrices are sparse, so the true gain is far smaller, but the structure is the same: blocks are factored independently, possibly in parallel, and the interface (a Schur complement or a set of boundary variables) couples them. A large interface erases the gain, so partitioners cut where coupling is weak, typically at MOSFET gates, which draw no DC current.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter on fast simulation; White and Sangiovanni-Vincentelli, "Relaxation Techniques for the Simulation of VLSI Circuits" (Kluwer, 1987).
