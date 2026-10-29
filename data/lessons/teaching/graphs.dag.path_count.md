# Counting paths through a node

The number of input-to-output paths in a netlist can be exponential in its size, which is why timers cannot enumerate paths and why path-based fault models and path coverage metrics need care. Counting paths is itself easy on a DAG, with the same dynamic program as longest paths but with sums in place of maxima. Path counts through a vertex also measure how much a single gate matters.

## Counting paths by dynamic programming

Let $N(v)$ be the number of paths from any source to $v$. Set $N(s) = 1$ at each source and $N(v) = \sum_{u \to v} N(u)$ otherwise, in topological order. The total number of source-to-sink paths is $\sum_{\text{sinks } t} N(t)$.

With sources $v_1, v_2$ and arcs $v_1 \to v_4$, $v_2 \to v_3$, $v_2 \to v_4$, $v_2 \to v_6$, $v_3 \to v_5$, $v_4 \to v_5$: $N(v_3) = 1$, $N(v_4) = 2$, $N(v_5) = 1 + 2 = 3$, $N(v_6) = 1$. The sinks are $v_5$ and $v_6$, so there are $3 + 1 = 4$ paths. The pass is linear, though the counts themselves can need many bits.

## Paths through one vertex

Run the same recurrence backward: $M(v) = \sum_{v \to w} M(w)$ with $M(t) = 1$ at sinks counts paths from $v$ to any sink. Every path into $v$ combines with every path out of it, so the number of source-to-sink paths through $v$ is $N(v) \cdot M(v)$.

If three paths reach a gate and four leave it, twelve paths pass through it. Path-based test generation, path delay fault coverage and criticality metrics all use these counts. A vertex with a huge product is a convergence point where many paths share hardware, so it is a natural target when an optimizer wants to improve many paths at once.

## Exponential path explosion

Reconvergent fanout multiplies paths. A chain of $k$ stages, each fanning out into $w_i$ parallel gates that reconverge, has $\prod_i w_i$ paths but only about $1 + \sum_i (w_i + 1)$ vertices. Ten stages of width 4 give $4^{10} \approx 10^6$ paths through about 50 vertices. Real datapaths, such as multipliers, are far worse.

This is the decisive argument for block-based timing: a per-vertex pass costs $O(V + E)$ while path enumeration costs time proportional to the number of paths. It is also why path-delay fault testing targets only selected paths, and why reports list the worst $k$ paths per endpoint using a best-first search over the DAG rather than all of them.

## Going deeper

Bushnell and Agrawal, "Essentials of Electronic Testing" (2000), chapter on path delay faults; Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 8.
