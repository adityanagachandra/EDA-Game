# The Fiduccia-Mattheyses algorithm

The Fiduccia-Mattheyses algorithm is the local-search engine at the heart of almost every netlist partitioner, from flat bipartitioning to the refinement step of hMETIS and KaHyPar. It runs passes of single-cell moves under a balance tolerance, locks each moved cell, and commits the best prefix of each pass. Understanding the pass structure is what lets a designer predict its behaviour and its limits.

## One pass of tentative moves

At the start of a pass all cells are free and their gains are in the buckets. The pass repeatedly picks the free cell with the highest gain whose move keeps balance, moves it, **locks** it so it cannot move again in this pass, updates its neighbours' gains, and records the gain. The pass ends when no free cell can move.

Locking guarantees termination and bounds a pass to at most $n$ moves. Moves are tentative: the pass may make many moves with zero or negative gain, because later moves may pay them back. The recorded gain sequence, for example $0, 2, 3, -2, 0, -2, -2$, is the raw material for the commit decision.

## Committing the best prefix

After the pass, FM computes cumulative gains $G_j = g_1 + \dots + g_j$ and keeps the prefix with the largest $G_j$, undoing every later move. For gains $0, 2, 3, -2, 0, -2, -2$ the cumulative gains are $0, 2, 5, 3, 3, 1, -1$, so the best prefix is the first three moves. Starting from a cut of 15, the committed cut is $15 - 5 = 10$.

Two wrong rules are tempting. Stopping at the first negative gain misses recoveries later in the pass. Keeping all moves can end worse than the start, here at 16. If the best cumulative gain is 0 (the empty prefix), the pass made no progress and the algorithm stops; otherwise a new pass starts from the committed partition.

## Strengths and limits

With gain buckets, a pass runs in $O(P)$ time, and only a few passes are usually needed, so FM is fast enough to run millions of times inside larger flows. It handles hyperedges and cell areas directly, which KL could not.

Its weakness is locality. Moving one cell of a tightly connected cluster rarely has positive gain, so FM seldom relocates a whole cluster, and on large netlists it converges to poor local minima that depend on the starting partition. Remedies are many random starts (expensive), smarter tie-breaking, lookahead gains (Krishnamurthy's higher-level gains), and above all the multilevel scheme, which runs FM on coarsened netlists where one move relocates a cluster.

## Going deeper

Fiduccia and Mattheyses (Proc. DAC, 1982); Krishnamurthy, "An improved min-cut algorithm for partitioning VLSI networks" (IEEE Trans. Computers, 1984); Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 2.
