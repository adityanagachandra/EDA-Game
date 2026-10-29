# R-trees for spatial queries

The R-tree is the standard dynamic index for rectangles: a height-balanced tree whose nodes store the minimum bounding rectangles (MBRs) of their children. Window queries descend only into children whose MBR meets the window. It handles insertions and deletions, works for shapes of any size, and is the index behind many layout databases and spatial databases. Its performance depends entirely on how much node MBRs overlap.

## Node bounding rectangles

A leaf stores entries (shape id plus bounding box); an internal node stores, for each child, the MBR of everything below it. The MBR of a set of boxes is $[\min x_1, \max x_2] \times [\min y_1, \max y_2]$. Every node holds between $m$ and $M$ entries (the root excepted), and all leaves sit at the same depth, so the height is $O(\log_m n)$. Example: entries $[3, 6] \times [4, 7]$, $[5, 9] \times [2, 5]$ and $[4, 8] \times [6, 10]$ give the MBR $[3, 9] \times [2, 10]$, area 48, although the entries cover far less. That empty area inside MBRs is the dead space that makes queries visit nodes containing nothing relevant.

## Descending into children

A window query starts at the root and descends into every child whose MBR intersects the window, pruning the others together with their whole subtrees. At the leaves it tests each entry and reports those that intersect. Example: a root with children MBRs $[2, 16] \times [3, 17]$, $[23, 37] \times [5, 19]$ and $[11, 25] \times [23, 37]$ and a window $[14, 22] \times [10, 20]$ descends only into the first child: the second starts at $x = 23$, beyond the window, and the third starts at $y = 23$. Unlike a k-d tree's cells, R-tree MBRs may overlap, so a query point can lie in several children and all of them must be searched; there is no worst-case logarithmic bound for window queries.

## Counting the work of a query

The cost of a query is the number of rectangle tests: every entry of every visited node. For a two-level tree that is the number of root entries plus the entries of each child that is visited. Example: a root with 3 children, of which 2 meet the window and hold 4 and 3 entries, costs $3 + 4 + 3 = 10$ tests, and perhaps 2 entries are reported. A brute-force scan of the 11 entries would cost 11 tests; the saving grows with depth and data size, since each level prunes most of the tree. Query cost rises with node overlap and dead space, which is why split policies and bulk loading, which control the shape of MBRs, decide R-tree quality.

## Going deeper

Guttman, "R-trees: a dynamic index structure for spatial searching" (ACM SIGMOD, 1984); Manolopoulos, Nanopoulos, Papadopoulos and Theodoridis, "R-Trees: Theory and Applications" (2006).
