# Balance constraints on a bipartition

Without a balance constraint, the best partition is trivial: put everything on one side and cut nothing. Every useful partitioning problem therefore bounds the size of each block. How the bound is stated, how a move is checked against it, and how large the resulting search space is are the practical details a partitioner designer must get right.

## The tolerance window

Let $A$ be the total cell area. A bipartition with imbalance tolerance $\epsilon$ requires each side's area $a_i$ to satisfy

$$(1 - \epsilon)\frac{A}{2} \le a_i \le (1 + \epsilon)\frac{A}{2}.$$

With $A = 286$ and $\epsilon = 0.2$, each side must hold between 114.4 and 171.6. Some tools state the bound only as an upper limit $a_i \le (1 + \epsilon) A / k$ for $k$ blocks, which implies a lower limit for bisection. Others count cells rather than area; for real netlists, with cells from tiny inverters to large flip-flops and macros, area is the meaningful measure.

## Checking a move

A move of a cell of area $s$ from side $A$ to side $B$ changes the side areas to $a_A - s$ and $a_B + s$. It is legal only if both remain inside the window. Large cells are therefore harder to move: with sides holding 35 and 32 and a window from 26.8 to 40.2, a cell of area 10 cannot move from $A$ to $B$ (B would reach 42), while a cell of area 3 can.

FM checks this for every candidate move and skips the best-gain cell when it is illegal. A tolerance that is too tight blocks many good moves and traps the search; one that is too loose produces lopsided blocks that hurt the next stage. That is why tolerances of a few percent are typical, with fixed macros handled before partitioning.

## Why exhaustive search is hopeless

For $n$ unit-area cells with $n$ even, the number of exact bisections into two unlabelled halves is $\binom{n}{n/2}/2$. For 8 cells that is 35; for 20 cells it is 92,378; for 100 cells it is about $5 \times 10^{28}$. By Stirling's formula the count grows like $2^n/\sqrt{n}$.

Minimum bisection is NP-hard, even for graphs, so no algorithm is known that avoids this growth in the worst case. Practical partitioners use local search (KL, FM) refined inside a multilevel hierarchy, and they are judged empirically on benchmark suites. Exact methods such as branch and bound or integer programming are limited to small instances, where they serve as references for testing heuristics.

## Going deeper

Garey and Johnson, "Computers and Intractability" (1979), on graph bisection; Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 2.
