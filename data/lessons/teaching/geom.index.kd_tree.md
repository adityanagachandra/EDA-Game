# k-d trees and nearest neighbors

A k-d tree splits space recursively, alternating axes, at the median of the points: x at the root, y at the next level, and so on. It answers range queries and nearest-neighbour queries over points, such as cell centres, pin locations or legal sites, by pruning cells that cannot contain an answer. It is static by nature but cheap to build in $O(n \log n)$ and very fast on point data.

## Splitting at the median

To build a 2-d tree, sort the points by x and split at the median: the median point is stored at the root, points with smaller x go left and larger x go right. Each child then splits its own points at their median y, and the axes alternate down the tree. Splitting at the median keeps the tree balanced with depth $\lceil \log_2 (n + 1) \rceil$. Example: points with x-coordinates $3, 8, 12, 17, 21, 26, 30$ split at $x = 17$, the fourth of seven. Presorting by both coordinates, or using linear-time selection, gives $O(n \log n)$ construction. The tree partitions the plane into rectangular cells, one per leaf, with no overlap, unlike an R-tree.

## Alternating axes

At depth one the splitting coordinate switches to y. Each child takes only its own points and splits them at their median y. Example: if the three points left of the root have y-values 5, 19, 33, the left child splits at $y = 19$; the right child splits its own points independently, at a different y. The resulting cells adapt to the data: dense regions get small cells and empty regions large ones. For range queries on $n$ points in two dimensions the k-d tree visits $O(\sqrt{n} + k)$ nodes in the worst case, since a query boundary line can cross only $O(\sqrt{n})$ cells. Alternating axes is the simplest rule; splitting along the axis of largest spread often works better on skewed data.

## Pruning a nearest-neighbour search

To find the nearest point to a query $q$, descend to the leaf cell containing $q$ and take the best candidate there, at distance $r$. Then unwind: a cell must be searched only if its distance from $q$ is less than $r$, because nothing inside it can be closer otherwise. In the Manhattan metric the distance from $q$ to a cell is its x-gap plus its y-gap, zero if $q$ is inside. Example: with $r = 6$, a cell 4 units to the right and 3 units above $q$ is at L1 distance 7 and is pruned, while a cell 2 to the right with an overlapping y-range is at distance 2 and is searched. Every improved candidate shrinks $r$ and prunes more, which gives about $O(\log n)$ expected time on well-distributed data.

## Going deeper

Bentley, "Multidimensional binary search trees used for associative searching" (Communications of the ACM, 1975); Friedman, Bentley and Finkel, "An algorithm for finding best matches in logarithmic expected time" (ACM Transactions on Mathematical Software, 1977); de Berg et al. (3rd ed., 2008), chapter 5.
