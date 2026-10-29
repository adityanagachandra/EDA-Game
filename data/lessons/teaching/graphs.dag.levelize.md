# Levelization of a combinational netlist

Levelization assigns each vertex of a combinational DAG a level: primary inputs and register outputs at level 0, and every other vertex one above its deepest fanin. Levels give a ready-made schedule for evaluation, a measure of logic depth, and the bucket key that incremental timers use to process dirty vertices in order. It is one of the first passes every simulator and timer runs.

## Levels from fanins

Define $\ell(v) = 0$ for a vertex with no incoming arcs, and $\ell(v) = 1 + \max_{u \to v} \ell(u)$ otherwise. This is the longest path, counted in arcs, from any source to $v$. It is computed in one pass over a topological order.

With arcs $v_1 \to v_4$, $v_2 \to v_4$, $v_4 \to v_5$, $v_4 \to v_6$, $v_3 \to v_6$, $v_5 \to v_7$: $v_1, v_2, v_3$ are at level 0, $v_4$ at level 1, $v_5$ and $v_6$ at level 2 (since $v_6$'s deepest fanin $v_4$ is at level 1), and $v_7$ at level 3. Every arc goes from a lower level to a strictly higher one, so the levels themselves are a valid evaluation schedule.

## Depth and width

The number of levels is the **logic depth** plus one, counting level 0. Depth correlates with delay and is what synthesis minimizes when it optimizes for speed before timing information is available. The number of vertices at a level is the **width** of that level, the parallelism available to a levelized simulator or a parallel timer that processes one level at a time.

In the example, the levels hold $\{v_1, v_2, v_3\}$, $\{v_4\}$, $\{v_5, v_6\}$ and $\{v_7\}$: four levels, widths 3, 1, 2, 1. A narrow level is a synchronization bottleneck for parallel evaluation. A long tail of narrow levels usually means a deep carry chain or a long serial path, which synthesis may restructure.

## Levels in incremental analysis

After an edit, an incremental timer must recompute arrival times only for dirty vertices, and it must process them in topological order. Keeping a level for every vertex makes this cheap: dirty vertices go into a priority queue or a bucket array keyed by level, and the timer pops them lowest level first. A vertex is then never recomputed before its fanins.

Edits can change levels. Inserting a buffer on an arc raises the level of everything downstream of that arc by up to one; deleting a gate can lower levels. Timers either propagate level changes eagerly or tolerate stale levels that remain topologically consistent, since correctness only needs every arc to point to a higher level, not the exact longest-path value. This distinction, an exact invariant versus a sufficient one, recurs throughout incremental engine design.

## Going deeper

Abramovici, Breuer and Friedman, "Digital Systems Testing and Testable Design" (1990), the chapter on logic simulation; Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 8.
