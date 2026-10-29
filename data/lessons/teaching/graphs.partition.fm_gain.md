# Fiduccia-Mattheyses cell gain

The Fiduccia-Mattheyses algorithm moves one cell at a time on the hypergraph itself, under a balance tolerance. Its gain is defined directly on nets: how many cut nets a move would uncut, minus how many uncut nets it would cut. This per-net view is what makes FM handle real netlists and, with the bucket structure, run in linear time per pass.

## Nets that a move uncuts

Let cell $c$ sit on its from-side $F$. A net containing $c$ becomes uncut by the move exactly when $c$ is the only pin of that net on $F$: all its other pins are already on the to-side $T$, and moving $c$ brings the whole net together. The number of such nets is $FS(c)$.

For nets $n_2 = \{c_1, c_3, c_4, c_5\}$ with only $c_3$ on side $B$, moving $c_3$ to $A$ uncuts $n_2$, so $n_2$ counts toward $FS(c_3)$. A net where $c$ shares its side with another pin does not count, because after the move that other pin still holds the net across the cut.

## Nets that a move cuts, and the gain

A net containing $c$ becomes cut by the move exactly when all of its pins are on $F$, including $c$: the move leaves $c$ alone on $T$. The number of such nets is $TE(c)$. The FM gain is

$$g(c) = FS(c) - TE(c),$$

the decrease in the cut. Every other net of $c$ is cut both before and after, so it does not change the gain. A cell with $FS = 1$ and $TE = 0$ has gain 1; a cell sitting in the middle of two uncut nets has $FS = 0$, $TE = 2$ and gain -2. Gains are integers between $-\deg(c)$ and $\deg(c)$ for unit-weight nets, and with weights each count becomes a sum of net weights.

## Choosing among cells

FM moves the free cell with the largest gain whose move keeps balance. Ignoring balance, the best cell in a small partition is found by evaluating $FS - TE$ for each cell; for example gains 2, 1, -1, -1, 0, 0, 1, 2 give a maximum of 2, shared by two cells. Ties are frequent, so the tie-breaking rule matters.

Computing all gains from scratch costs $O(P)$, one look at every pin. FM's insight is that after a move only cells on a few nets change gain, so gains are maintained incrementally in buckets rather than recomputed. The gain is also local: it sees only the next move, so a sequence of individually poor moves can be needed to escape a local minimum, which FM handles by continuing the pass and keeping the best prefix.

## Going deeper

Fiduccia and Mattheyses, "A linear-time heuristic for improving network partitions" (Proc. Design Automation Conference, 1982); Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 2.
