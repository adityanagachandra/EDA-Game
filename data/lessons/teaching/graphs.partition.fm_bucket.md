# Gain buckets and incremental gain updates

What makes FM linear per pass is not the gain formula but its data structure: an array of buckets indexed by gain, each a doubly linked list of cells, with a pointer to the highest non-empty bucket. Combined with update rules that touch only a few nets after each move, it gives constant-time selection and update. An engine designer must know its size, how updates propagate, and where it breaks down.

## Sizing the bucket array

For unit-weight nets, a cell of degree $p$ has gain between $-p$ and $p$. So each side's bucket array needs $2p_{\max} + 1$ buckets, where $p_{\max}$ is the largest cell degree in the netlist. With degrees 3, 4, 3, 4, 2, 2, 4, $p_{\max} = 4$ and each array has 9 buckets, for gains -4 through 4.

Each bucket holds the free cells with that gain in a doubly linked list, and an index into the array tracks the maximum non-empty bucket. Selecting the best cell is reading that bucket; when it empties, the pointer moves down. With weighted nets the gain range becomes the largest weighted degree, which can be huge; tools then quantize weights or switch to a heap, losing strict linear time.

## Updating gains after a move

When cell $c$ moves from $F$ to $T$, only nets of $c$ can change any gain, and only when they are **critical**: 0 or 1 pins on $T$ before the move, or 0 or 1 pins on $F$ after it. Before the move: if $T$ held 0 pins, the net was uncut and sat in the $TE$ count of every other cell on it, so each free cell on the net gains 1; if $T$ held exactly 1 pin, that cell loses the net from its $FS$, so its gain drops by 1. After the move: if $F$ holds 0 pins, every free cell on the net loses 1; if $F$ holds exactly 1 pin, that cell gains 1.

The safe way to check an update rule is recomputation: compute every cell's $FS - TE$ on the old and new partitions and compare. For a cell whose only shared net with $c$ goes from having one pin on its side to two, the gain changes; for a cell sharing only non-critical nets, it does not.

## Why updates touch few cells

Each update is a constant-time unlink and relink between buckets. A net's state changes only a bounded number of times per pass once its pins are locked, so the total update work over a pass is $O(P)$. In practice, after one move only a handful of cells change gain: those on the moved cell's critical nets.

Two engineering consequences follow. First, the order cells are inserted into a bucket (LIFO or FIFO) breaks ties and measurably changes results; LIFO tends to keep moving cells near recent moves, which helps. Second, the incremental bookkeeping is easy to get wrong in edge cases (nets with duplicate pins, nets entirely of locked cells), so engines verify the incrementally maintained cut against a full recount in testing.

## Going deeper

Fiduccia and Mattheyses, "A linear-time heuristic for improving network partitions" (Proc. DAC, 1982); Hagen, Huang and Kahng, "On implementation choices for iterative improvement partitioning algorithms" (IEEE TCAD, 1997).
