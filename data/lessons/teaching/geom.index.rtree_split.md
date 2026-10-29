# R-tree node splitting

When an R-tree node overflows, its entries must be divided between two nodes, and the division decides how much the resulting MBRs overlap and how much dead space they enclose. Every later query pays for a poor split. Guttman's quadratic split and the R*-tree's criteria of overlap and margin are the classic policies, and both are simple area and perimeter computations on boxes.

## Dead area of a pair

Guttman's quadratic split measures how badly two entries fit together by the dead area they would waste in a shared node:

$$d(E_1, E_2) = \text{area}(\text{MBR}(E_1 \cup E_2)) - \text{area}(E_1) - \text{area}(E_2).$$

Example: $E_1 = [0, 3] \times [0, 2]$ and $E_2 = [8, 10] \times [5, 9]$ have a joint MBR $[0, 10] \times [0, 9]$ of area 90, so $d = 90 - 6 - 8 = 76$, a lot of empty space. Two nearby boxes give a small $d$; overlapping boxes can make $d$ negative, since their overlap is subtracted twice. The quantity is cheap to compute and captures the intuition that entries far apart should not share a node. Because it depends only on the two boxes, it can be computed for every pair before any assignment decision is made.

## Picking seeds

PickSeeds evaluates $d$ for every pair of the $M + 1$ entries and takes the pair with the largest dead area as the seeds of the two new nodes: they are the worst pair to keep together. The remaining entries are then assigned one at a time, each to the node whose MBR it enlarges least, choosing first the entry with the strongest preference. The seed step costs $O(M^2)$, hence the name quadratic; a linear variant picks seeds by the largest normalized separation along each axis. Example: among four entries, scoring the six pairs and finding the maximum $d = 76$ for the far-apart pair makes those two the seeds, and the other entries join whichever seed is closer in enlargement.

## Overlap and margin in the R*-tree

The R*-tree evaluates candidate splits more carefully. For each axis it sorts entries by their lower and upper coordinates and considers every split into two groups of allowed sizes. It picks the split axis by the smallest sum of margins (the perimeters of the two group MBRs), which favours square-ish nodes, and then the split by the smallest overlap area between the two group MBRs, breaking ties by total area. Example: groups with MBRs $[2, 9] \times [1, 6]$ and $[7, 15] \times [4, 12]$ have margins 24 and 32, sum 56, and overlap $2 \times 2 = 4$. Overlap matters most: a query falling in the overlap must descend into both nodes. The R*-tree also reinserts some entries on overflow, which improves the tree further.

## Going deeper

Guttman (1984), as above; Beckmann, Kriegel, Schneider and Seeger, "The R*-tree: an efficient and robust access method for points and rectangles" (ACM SIGMOD, 1990).
