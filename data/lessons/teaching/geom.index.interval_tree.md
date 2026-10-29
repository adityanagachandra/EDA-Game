# Interval trees

An interval tree stores $n$ intervals and answers stabbing queries, all intervals containing a point, in $O(\log n + k)$ time with $O(n)$ storage. It is the one-dimensional workhorse behind sweep statuses, cut-line queries in routers, and the y-dimension of rectangle overlap reporting. The centred version is short to describe and shows clearly why the query never wastes work.

## Stabbing queries

Given closed intervals $[a_i, b_i]$ and a query $q$, the stabbing query reports every interval with $a_i \le q \le b_i$. In layout this is: which wires cross the cut line $y = q$, which rows a coordinate falls in, which shapes are active at a sweep position. Example: intervals $[2, 9]$, $[5, 14]$, $[10, 18]$, $[12, 20]$ and query 12 report $[5, 14]$, $[10, 18]$, $[12, 20]$, three intervals, since 12 lies in $[12, 20]$ as a closed endpoint. A linear scan costs $O(n)$ per query; an interval tree pays $O(n \log n)$ once to build and then answers each query in time proportional to $\log n$ plus the answer size, which matters when there are millions of queries. The query is also the core operation of a sweep status, where the stabbing point is the current sweep position.

## Building a centred interval tree

Choose a centre $x_c$, typically the median of all endpoints. Intervals containing $x_c$ are stored at the root; intervals entirely left of $x_c$ go to the left subtree and those entirely right go to the right subtree, recursively. The root keeps its intervals in two lists, sorted by left endpoint ascending and by right endpoint descending. Because the median splits the endpoints evenly, the depth is $O(\log n)$, and every interval is stored exactly once, so storage is $O(n)$. Example: with centre 15, $[10, 18]$ and $[12, 20]$ stay at the root, $[2, 9]$ goes left, and $[16, 25]$ goes right. Intervals are never split, unlike in a segment tree.

## Querying the root

Every interval at the root contains $x_c$. If the query $q$ is left of $x_c$, a root interval contains $q$ exactly when its left endpoint is at most $q$, because its right endpoint is already beyond $x_c > q$. So scan the list sorted by left endpoint and stop at the first left endpoint greater than $q$. Symmetrically, for $q > x_c$ scan by right endpoint descending and stop at the first right endpoint less than $q$. Then recurse into one child only. Example: root centre 15 with $[10, 18]$, $[12, 20]$, $[14, 16]$ and $q = 13$: left endpoints 10, 12, 14, so two are reported and the scan stops at 14. Each node costs $O(1 + k_{\text{node}})$, and the path has $O(\log n)$ nodes.

## Going deeper

Edelsbrunner, "Dynamic data structures for orthogonal intersection queries" (Technical University of Graz report, 1980); de Berg et al., "Computational Geometry: Algorithms and Applications" (3rd ed., 2008), chapter 10.
