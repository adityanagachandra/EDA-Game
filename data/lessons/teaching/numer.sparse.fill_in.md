# Fill-in and reordering

Gaussian elimination on a sparse matrix creates new nonzeros, called fill, in positions that were zero. Fill decides the memory and time of every direct solve in a simulator, and it depends entirely on the elimination order. Seeing elimination as a game on the matrix graph turns the order question into graph reasoning, which is how every ordering algorithm works.

## Elimination as a graph game

For a structurally symmetric matrix, draw a vertex per unknown and an edge per off-diagonal nonzero pair. Eliminating unknown $v$ updates every entry $a_{ij}$ with $i, j$ neighbours of $v$: $a_{ij} \leftarrow a_{ij} - a_{iv}a_{vj}/a_{vv}$. If $a_{ij}$ was zero it becomes nonzero (barring exact cancellation, which symbolic analysis ignores).

In the graph: remove $v$ and join all its remaining neighbours into a clique. Each new edge is a fill pair, one new nonzero in $L$ below the diagonal and its mirror in $U$. Example: vertex 3 with neighbours 1, 4 and 6, where only 1-4 is already an edge. Eliminating 3 adds edges 1-6 and 4-6: two fill edges.

## Counting fill for an elimination order

Repeat the game in the chosen order, always on the current graph, which includes earlier fill. Take a 5-cycle 1-2-3-4-5-1. Order 1, 2, 3, 4, 5: eliminating 1 joins 2 and 5 (one fill edge). Now 2 has neighbours 3 and 5, not joined: eliminating 2 adds 3-5 (fill 2). Then 3 has neighbours 4 and 5, already joined. Total fill: 2.

The symbolic phase of a sparse direct solver does exactly this bookkeeping, more cleverly, before any numbers are touched. It returns the pattern of $L$ and $U$, so memory can be allocated once. Total stored entries are $\text{nnz}(A) + 2 \cdot \text{fill}$ for a symmetric pattern with the diagonal counted once.

## Why order changes fill: the arrow matrix

Take a hub vertex joined to $k$ others that are not joined to each other, the pattern of a supply node or a bus. Eliminating the hub first joins all $k$ neighbours into a clique: $k(k - 1)/2$ fill edges, and the rest of the matrix is dense. Eliminating the hub last creates no fill at all, because each leaf has only the hub as a neighbour when it is removed.

For $k = 1000$ that is the difference between half a million fill entries and none. Finding a minimum-fill order is NP-hard (Yannakakis, 1981), so tools use heuristics: minimum degree and its approximate variant AMD, Markowitz pivoting for unsymmetric matrices, and nested dissection for mesh-like grids. A good order is computed once per pattern and reused across every Newton iteration and timestep.

## Going deeper

George and Liu, "Computer Solution of Large Sparse Positive Definite Systems" (Prentice-Hall, 1981); Davis, "Direct Methods for Sparse Linear Systems" (SIAM, 2006), chapter 7.
