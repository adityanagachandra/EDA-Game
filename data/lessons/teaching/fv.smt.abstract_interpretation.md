# Abstract interpretation with intervals

Abstract interpretation computes sound over-approximations of all the values a program or design can produce, by running it over an abstract domain instead of concrete values. The interval domain, which tracks a lower and upper bound for each variable, is the simplest useful example. Compilers, static analyzers and HDL tools use it for range analysis, bit-width inference and dead-code detection, and the same lattice-and-fixpoint machinery underlies model checking.

## Interval arithmetic

Each variable is abstracted by an interval $[l, h]$ containing all its possible values, and each operation is replaced by a sound interval version. With $x \in [1, 3]$ and $y \in [-2, 4]$: $x + y \in [1 - 2, 3 + 4] = [-1, 7]$; $x - y \in [1 - 4, 3 + 2] = [-3, 5]$; and $x \cdot y$ takes the minimum and maximum of the four corner products $-2, 4, -6, 12$, so $[-6, 12]$. The results are sound, containing every concrete result, but can be imprecise. The classic loss is $x - x$, which interval arithmetic evaluates as $[1 - 3, 3 - 1] = [-2, 2]$ instead of $[0, 0]$, because the domain is non-relational: it forgets that both operands are the same variable. Relational domains such as octagons or polyhedra track constraints between variables and avoid this loss at higher cost.

## Joins at merge points

Where control flow merges, after an if-else or at a loop head, the analysis joins the intervals from each incoming path: $[a, b] \sqcup [c, d] = [\min(a, c), \max(b, d)]$, the smallest interval containing both. With $x \in [2, 5]$ and an unknown condition, if (c) x = x + 4; else x = x * 2 gives $[6, 9]$ on one side and $[4, 10]$ on the other, and the join is $[4, 10]$. The join loses information again: the analysis cannot express that $x$ is never 5, since intervals have no holes. The intervals ordered by inclusion form a lattice with this join, and every abstract operation is monotone, which is what guarantees that the analysis computes a sound fixpoint. Hardware tools use exactly this computation to infer that a signal never exceeds a bound and can be narrowed to fewer bits.

## Loops, fixpoints and widening

At a loop head the interval is the least fixpoint of: join of the entry value and the value coming around the loop. For i = 0; while (i < N) i = i + s, the body sees the head interval intersected with the guard $i \le N - 1$, and adds $s$. With $N = 10$ and $s = 3$, Kleene iteration gives head intervals $[0, 0]$, $[0, 3]$, $[0, 6]$, $[0, 9]$, then the clipped body $[0, 9] + 3 = [3, 12]$ makes the head $[0, 12]$, which is stable. In general the head stabilizes at $[0, N - 1 + s]$ and the exit interval, the head intersected with $i \ge N$, is $[N, N - 1 + s]$, here $[10, 12]$. Intervals cannot record that $i$ is a multiple of 3. For large $N$ this iteration takes $N / s$ steps, so analyzers apply widening, jumping an unstable bound to infinity to force termination, and then narrowing, which iterates downward to recover a finite bound such as 12.

## Going deeper

Cousot and Cousot, "Abstract Interpretation: A Unified Lattice Model for Static Analysis of Programs by Construction or Approximation of Fixpoints" (POPL 1977); Rival and Yi, "Introduction to Static Analysis: An Abstract Interpretation Perspective" (2020).
