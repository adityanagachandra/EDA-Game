# Segment trees and canonical decomposition

A segment tree stores intervals by splitting each into $O(\log n)$ canonical pieces that are nodes of a balanced tree over the elementary intervals. It uses $O(n \log n)$ storage but supports what interval trees cannot easily do: aggregates over coverage, such as the covered length used by the union-area sweep, and multi-level structures for rectangle queries. Understanding canonical decomposition is the key to using it.

## Elementary intervals

Sort the distinct endpoints $p_1 < \dots < p_m$ of all intervals. They cut the line into elementary intervals: the open gaps $(p_i, p_{i+1})$, the single points $[p_i, p_i]$ and the two unbounded ends, $2m + 1$ pieces in all. Within one elementary interval every input interval either covers all of it or none of it, so the elementary intervals are the leaves of the tree. Example: intervals $[2, 7]$, $[4, 9]$, $[7, 12]$ have distinct endpoints 2, 4, 7, 9, 12, so $m = 5$ and there are 11 leaves. In integer layout coordinates many tools simplify to half-open unit cells $[i, i + 1)$ between compressed coordinates, which removes the point leaves. The tree over these leaves is built bottom-up in linear time once the endpoints are sorted.

## Canonical nodes

Each tree node covers the union of its leaves. An interval is stored at its canonical nodes: the highest nodes whose range lies entirely inside the interval and whose parent's range does not. Descend from the root; at a node fully inside the interval, store it and stop; at a node disjoint from it, stop; otherwise recurse into both children. At most two nodes per level are stored, so each interval uses $O(\log n)$ nodes. Example: over 8 unit leaves $[0, 8)$, the interval $[1, 7)$ decomposes into $[1, 2)$, $[2, 4)$, $[4, 6)$ and $[6, 7)$, four canonical nodes. A stabbing query walks one root-to-leaf path and reports every interval stored on it, in $O(\log n + k)$.

## Storage and aggregates

Summing canonical nodes over all intervals gives the storage, $O(n \log n)$ in the worst case, against $O(n)$ for an interval tree. The payoff is that each node's stored set is a clean unit for aggregates. The union-area sweep stores at each node a count of intervals covering it fully and the covered length below it: if the count is positive the covered length is the node's whole range, otherwise the sum of its children's. An insertion or deletion updates only the $O(\log n)$ canonical nodes and their ancestors. Example: over 16 leaves, intervals $[0, 16)$, $[3, 5)$ and $[2, 11)$ need 1, 2 and 4 entries, 7 in all. Multi-level segment trees, with a second structure at each node, answer rectangle queries in $O(\log^2 n + k)$.

## Going deeper

de Berg et al., "Computational Geometry: Algorithms and Applications" (3rd ed., 2008), chapter 10; Bentley, "Algorithms for Klee's rectangle problems" (1977).
