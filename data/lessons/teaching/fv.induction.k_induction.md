# k-induction

k-induction strengthens ordinary induction by assuming the property on $k$ consecutive states instead of one. Unreachable states that only lead out of the property through short paths stop being counterexamples to induction once $k$ is large enough. It is the proof method behind the prove mode of SymbiYosys's smtbmc engine and of many commercial engines, so knowing what $k$ a design needs, and when no $k$ works, is practical knowledge.

## Base case and step case

For depth $k$ there are two checks. Base case: no path of fewer than $k$ transitions from an initial state reaches a $\neg P$-state, which is BMC to depth $k - 1$. Step case: for every path $s_0, s_1, \dots, s_k$ of $k$ transitions, if $s_0$ to $s_{k-1}$ all satisfy $P$ then $s_k$ satisfies $P$. With $k = 1$ this is ordinary induction. If $P$ holds in every reachable state, every base case passes, and the step case fails exactly when some path of $k$ consecutive $P$-states is followed by a $\neg P$-state. Example: s0 ↔ s1 reachable, plus an unreachable chain s2 → s3 → s4 with $P = \{s0, s1, s2, s3\}$. For $k = 1$, s3 → s4 fails; for $k = 2$, the path s2, s3, s4 fails; for $k = 3$ no path of three $P$-states ends at s3, so $P$ is proved with $k = 3$.

## Computing the minimal k

When $P$ is an invariant, the minimal $k$ follows from the graph of $P$-states. Mark every $P$-state that has a $\neg P$ successor. For each, find the longest path of $P$-states ending at it, counting states. If that longest path has $L$ states, the step case fails for every $k \le L$ and holds for $k = L + 1$. If no $P$-state has a $\neg P$ successor, $P$ is already inductive and $k = 1$. Extend the example with s5 → s2, s6 → s5 and s7 → s3, all in $P$. The only $P$-state with a $\neg P$ successor is still s3. The longest $P$-path ending at s3 is s6, s5, s2, s3, with 4 states, so the minimal $k$ is 5. Engines simply try $k = 1, 2, 3, \dots$, with the base cases supplied by an incremental BMC run.

## Cycles and the simple-path constraint

If a cycle of unreachable $P$-states can reach a $P$-state with a $\neg P$ successor, the path of $P$-states can be made arbitrarily long by going around the cycle, so the step case fails for every $k$: plain k-induction never proves $P$. Replace the chain's s2 by a self-loop s2 → s2, keeping s2 → s3 → s4: the path s2, s2, ..., s2, s3, s4 defeats every $k$. The fix is the simple-path constraint: in the step case require $s_0, \dots, s_k$ to be pairwise distinct. A path that repeats a state can be shortened, so the constraint is sound, and with it some $k$ always works on a finite system. Then the minimal $k$ is 1 plus the number of states on the longest loop-free $P$-path ending at a marked state: here s2, s3, so $k = 3$. The price is a quadratic number of disequality constraints, and in the worst case $k$ is as large as the recurrence diameter.

## Going deeper

Sheeran, Singh and Staalmarck, "Checking Safety Properties Using Induction and a SAT-Solver" (FMCAD 2000); Een and Sorensson, "Temporal Induction by Incremental SAT Solving" (BMC Workshop, 2003).
