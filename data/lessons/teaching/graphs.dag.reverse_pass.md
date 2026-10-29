# Reverse-topological passes and node slack

A forward pass computes the latest arrival at every vertex. A backward pass, in reverse topological order, computes the latest time each vertex may settle without violating any endpoint requirement. Their difference is slack, the per-vertex margin that drives every timing optimizer. As a graph algorithm the backward pass is the min-plus mirror of the forward pass, and together they identify the critical subgraph.

## Required times flow backward

Give every endpoint a required time $q(s)$. For every other vertex,

$$q(v) = \min_{v \to w} \big(q(w) - d(v,w)\big),$$

computed in reverse topological order so that every successor is final. The minimum appears because the vertex must satisfy its tightest successor.

With arcs $v_1 \to v_4$ (1), $v_2 \to v_3$ (5), $v_2 \to v_4$ (5), $v_3 \to v_5$ (9), $v_4 \to v_5$ (7) and endpoint $v_5$ required at 34 ps: $q(v_3) = 25$, $q(v_4) = 27$, $q(v_2) = \min(25 - 5, 27 - 5) = 20$, $q(v_1) = 26$. The pass costs $O(V + E)$, the same as the forward pass, and it uses the same topological order read backward.

## Slack at a vertex

The slack of a vertex is $s(v) = q(v) - a(v)$: how much later the signal at $v$ could settle before some endpoint fails. Negative slack means at least one path through $v$ violates its requirement. In the example, arrivals are $a(v_1) = 0$, $a(v_4) = 5$, so $s(v_1) = 26$ and $s(v_4) = 22$.

Slack is a path property stored on vertices: $s(v)$ equals the slack of the worst endpoint path through $v$. So an optimizer reading a vertex's slack knows the worst case through that vertex without enumerating paths. Arc slack, $q(w) - a(v) - d(v,w)$, is the analogous quantity for an arc and tells a sizer which fanin of a gate is the critical one.

## The critical subgraph

When the required time is set to the latest endpoint arrival, the worst slack in the graph is zero, and the vertices with zero slack are exactly those on at least one longest path. They form the **critical subgraph**. Its size tells an optimizer how concentrated the problem is: one thin path can be fixed by sizing a few gates, while a wide critical region needs restructuring.

Both passes are linear, so a timer can recompute all slacks after each optimization step. The same pair of passes computes ASAP and ALAP times in high-level synthesis scheduling, where the difference is an operation's mobility, and early and late times in project scheduling. In every case, zero-slack vertices are the ones that cannot move without lengthening the whole schedule.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 8; De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 5 on ASAP and ALAP scheduling.
