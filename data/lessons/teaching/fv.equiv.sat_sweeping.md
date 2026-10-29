# SAT sweeping and simulation signatures

SAT sweeping, also called fraiging in ABC, finds and merges functionally equivalent nodes inside one AIG. It combines cheap random simulation, which rules out most false candidates, with SAT calls that prove the rest. In equivalence checking it collapses two similar designs into one graph from the inputs upward; in synthesis it removes redundancy. Its cost model, simulation classes and one SAT call per candidate, is easy to compute and drives engine tuning.

## Simulation signatures

Apply $k$ random input vectors (in practice as 64-bit words, so 64 vectors per pass) and record each node's output on each vector: its $k$-bit signature. If two nodes have different signatures, some vector distinguishes them, so they are certainly not equivalent. If their signatures are equal, they are candidates. Grouping nodes by signature gives candidate equivalence classes; a class with one member needs no further work. With signatures $n_1$: 1010, $n_2$: 0110, $n_3$: 1010, $n_4$: 0110, $n_5$: 1111, there are two candidate classes, $\{n_1, n_3\}$ and $\{n_2, n_4\}$, and $n_5$ is a singleton. More vectors split more false classes; the right number is a trade-off between simulation time and wasted SAT calls on candidates that fail.

## Proving candidates with SAT

For each class, pick a representative, typically the node closest to the inputs, and prove every other member equal to it with one SAT call on a small miter of the two nodes. A class of size $s$ costs $s - 1$ calls if every candidate holds; in the example, $1 + 1 = 2$ calls. Sweeping proceeds in topological order, and each proved pair is merged immediately, so that later queries see the merged graph and stay small; proved nodes can also be treated as cut points. A failing call returns a counterexample vector, which is not a bug in the design: it is added to the simulation patterns and re-run, which typically splits several other wrong classes at once. Engines also limit the conflicts per call and skip candidates that are too hard, since an unproved candidate only costs optimization, not soundness.

## Equivalence up to complement

Logic often contains a node and the complement of another, for example after De Morgan rewriting, and AIGs represent inverters on edges for free. So sweepers treat $u$ and $\neg v$ as candidates too: a signature and its bitwise complement belong in the same class. The usual trick is to normalize each signature, say by complementing it when its first (or last) bit is 1, and then group the normalized signatures. With $n_1$: 1100, $n_2$: 0011, $n_3$: 1100 and $n_4$: 1010, normalizing so the first bit is 0 maps $n_1$ and $n_3$ to 0011, the same as $n_2$, so $\{n_1, n_2, n_3\}$ is one class with mixed phases and $n_4$ is alone. Each proof then checks either $u = v$ or $u = \neg v$ according to the phases, and the merge connects readers through an inverted edge.

## Going deeper

Mishchenko, Chatterjee, Jiang and Brayton, "FRAIGs: A Unifying Representation for Logic Synthesis and Verification" (ERL technical report, 2005); Kuehlmann et al., "Robust Boolean Reasoning for Equivalence Checking and Functional Property Verification" (IEEE TCAD, 2002).
