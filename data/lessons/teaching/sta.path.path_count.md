# Counting paths in a timing graph

Why can a timer not simply list every path and time each one? Because the number of paths grows exponentially with depth while the graph grows only linearly. Counting paths takes one dynamic-programming pass, and the count shows at once why block-based propagation, which never lists paths, is the only way to time a full chip, and why path-based analysis must be selective.

## Counting with one forward pass

Let $P(v)$ be the number of distinct paths from any source to $v$. A source has $P = 1$, and every other node sums its fanins: $P(v) = \sum_{u \to v} P(u)$. Visit the nodes in topological order. For arcs $D \to N$, $D \to T$, $T \to N$, $N \to E$, $Q \to M$ and $M \to E$, the counts are $D = 1$, $T = 1$, $N = 2$, $Q = 1$, $M = 1$ and $E = 3$. This is the same recurrence as arrival propagation with $(\max, +)$ replaced by $(+, \times)$: both are path problems over a semiring, which is why the same levelized traversal solves them.

## Paths through a node

Every path through $v$ is a path into $v$ joined to a path out of $v$, so the count through $v$ is $P_{in}(v) \times P_{out}(v)$. $P_{out}$ comes from the mirror pass: a sink has one path out, and every other node sums its fanouts. If four paths reach a node and two leave it, eight source-to-sink paths pass through it. The same split of in and out quantities gives node slack: the latest arrival in and the tightest required time out. A timer can therefore rank pins by how many paths they carry, which is useful for deciding where an optimization or a buffering change has the broadest effect.

## Exponential growth with reconvergence

Cascade $n$ stages, each fanning out into $w$ parallel branches that rejoin: the path count is $w^n$, while the graph has only about $(w+1)n$ arcs. Stages with 2, 2, 3 and 3 branches give $2 \times 2 \times 3 \times 3 = 36$ paths. Real datapaths, with adders, multipliers and multiplexers, reconverge constantly, and a modest block easily has more paths than could ever be listed. This is why path-based analysis re-times only the worst paths, found by enumeration in slack order, and why any algorithm that is linear in the number of paths is unusable on its own. False-path detection over all paths is hard for the same reason.

## Going deeper

Sapatnekar, "Timing" (2004), the chapter on static timing analysis; Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms", dynamic programming on DAGs.
