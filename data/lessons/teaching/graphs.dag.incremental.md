# Incremental updates on a DAG

Optimization loops change one gate, one wire or one buffer at a time and need updated timing after each change. Recomputing every arrival from scratch is linear in the design but far too slow when repeated a million times. Incremental algorithms recompute only the region a change can affect, and stop early where the change is absorbed. Designing that bookkeeping correctly is one of the central engineering problems of timers.

## The fanout cone bounds the work

If the delay of arc $u \to v$ changes, only arrivals at vertices reachable from $v$ can change, including $v$ itself. That set is the fanout cone of $v$, and an incremental timer marks it dirty. Required times change only in the fan-in cone of the changed arc's tail, $u$.

In a graph with arcs $v_1 \to v_3$, $v_3 \to v_4$, $v_3 \to v_7$, $v_4 \to v_7$, $v_1 \to v_5$, a change on $v_1 \to v_3$ marks $\{v_3, v_4, v_7\}$: three vertices, while $v_5$ is untouched. For a change near the outputs the cone is tiny; for a change near the inputs it can be most of the design. Timers keep the dirty set in a level-keyed queue so each vertex is recomputed after its fanins.

## Stopping early

Not every vertex in the cone actually changes. If arc $v_2 \to v_5$ slows by 5 ps but $v_5$'s arrival is set by another fanin with more margin than that, $a(v_5)$ is unchanged, and nothing downstream of $v_5$ needs recomputation through that route. An incremental timer recomputes a dirty vertex, compares the new arrival with the old, and propagates dirtiness to its successors only if the value changed.

So the work is proportional to the vertices whose values change plus their immediate fanouts, which is often far smaller than the cone. The test must be exact (or a tolerance that the signoff flow accepts), because skipping a vertex whose value changed by a tiny amount leaves stale values that later comparisons trust.

## Correctness of incremental results

An incremental engine is correct only if its results equal a full recomputation. That makes full recomputation the natural differential test: after a random sequence of edits, rerun the whole graph and compare every arrival and required time. Typical bugs are a missed dirty mark (an edit path that forgets to notify), stale levels after structural edits, and early stopping on a value that changed in a field not compared, such as slew.

Structural edits are harder than delay edits: inserting a buffer adds vertices and arcs, so levels and topological positions must be updated before propagation. Many timers handle structural edits by repairing levels locally and then running the same dirty-cone propagation. The speed of this loop is what makes timing-driven placement, sizing and ECO flows practical.

## Going deeper

The OpenSTA and OpenTimer documentation on incremental timing; Huang, Guo, Lin and Wong, "OpenTimer v2" (IEEE TCAD, 2021).
