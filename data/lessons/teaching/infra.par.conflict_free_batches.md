# Conflict-free parallel batches

Many EDA engines work on a grid of regions: detailed routers route tiles, legalizers repair windows, DRC engines check windows with halos. Neighbouring regions can touch the same wires or cells, so processing them simultaneously would race. Locks per object are slow and nondeterministic. A cleaner scheme colors the regions so that no two neighbours share a color, then processes one color at a time: inside a color, regions are independent and run in parallel with no locks at all.

## Coloring tiles so batches never conflict

Model each tile as a vertex and join two tiles that may conflict; a batch must be an independent set, and a coloring partitions the tiles into batches. For a rectangular grid where tiles conflict when they share an edge or a corner, the coloring by $(r \bmod 2,\ c \bmod 2)$ works: two tiles with the same color differ by at least 2 in row or column, so they never touch. Four colors are also necessary, since any 2 × 2 block contains four mutually adjacent tiles.

Example: a 5 × 7 grid. Rows 0, 2 and 4 are even (3 rows); columns 0, 2, 4 and 6 are even (4 columns). The even-even batch has $3 \times 4 = 12$ tiles. The other batches have $3 \times 3 = 9$, $2 \times 4 = 8$ and $2 \times 3 = 6$, which sum to the 35 tiles. If only edge-sharing tiles conflicted, a checkerboard with two colors would do.

## Rounds with a fixed worker count

Batches run one after another, separated by a barrier, because a tile in a later batch may touch tiles of the current one. Inside a batch of $n_k$ tiles, $p$ workers handle $p$ tiles per round, so the batch takes $\lceil n_k / p \rceil$ rounds, and the whole grid takes
$$R = \sum_{k} \left\lceil \frac{n_k}{p} \right\rceil.$$
For the 5 × 7 grid with $p = 4$: $\lceil 12/4 \rceil + \lceil 9/4 \rceil + \lceil 8/4 \rceil + \lceil 6/4 \rceil = 3 + 3 + 2 + 2 = 10$ rounds. The ceilings are the cost of the barriers: a batch of 9 tiles leaves three workers idle in its last round. Engines reduce that waste with more tiles than workers, so the last partial round is a small part of each batch, or by letting idle workers begin the next batch's tiles that do not touch unfinished ones.

## Speedup of a batched schedule

If every tile takes the same time $t$, the serial run takes $nt$ for $n$ tiles and the batched run takes $Rt$, so the speedup is $n/R$, at most $p$. For the example, $35/10 = 3.5$ on 4 workers: the partial final rounds cost 12.5%.

The scheme has two further virtues. It is deterministic, because which tiles run together and in what batch order is fixed by the coloring, not by timing, so the results do not depend on thread scheduling. And it needs no locks inside a batch. The trade-off is idle time at barriers and sensitivity to uneven tiles: one congested tile makes its whole batch wait. Real routers therefore size tiles to balance work and process batches with a task pool, but keep the coloring that makes concurrent tiles independent.

## Going deeper

The TritonRoute papers on parallel detailed routing (Kahng, Wang and Xu); Jones and Plassmann, "A Parallel Graph Coloring Heuristic" (SIAM J. Scientific Computing, 1993).
