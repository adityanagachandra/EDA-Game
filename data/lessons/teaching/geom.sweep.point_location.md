# Point location by slabs and trapezoidal maps

Point location asks which region of a planar subdivision contains a query point: which polygon a probe lands on, which net a via hits, which tile of a partition a coordinate falls in. With preprocessing, queries take $O(\log n)$. The slab method shows the idea simply but can need quadratic storage; the trapezoidal map achieves linear size and logarithmic expected query time and is the structure behind many practical locators.

## Slabs from vertical lines

Draw a vertical line through every segment endpoint. Within each slab between consecutive lines, no segment starts or ends, so the segments crossing the slab are totally ordered from bottom to top. With $m$ distinct endpoint x-values there are $m + 1$ slabs, counting the two unbounded ones. A query does two binary searches: one over the slab boundaries to find the slab, one over the slab's ordered segments using orientation tests to find the segment just below the point. Example: segments with endpoints at $x = 2, 5, 5, 9, 11, 14$ have five distinct x-values and six slabs. Both searches are $O(\log n)$, so queries are fast, and the structure is simple to build with a sweep.

## The cost of slabs

Each slab stores its own sorted list of the segments spanning it. A long segment spans many slabs and is stored in every one. Total storage is the sum over slabs of the number of spanning segments, which can reach $\Theta(n^2)$: picture $n/2$ long horizontal segments stacked over $n/2$ short ones that create $n$ slab boundaries. Example: with boundaries at $x = 0, 4, 6, 10$ and segments covering $[0, 10]$, $[0, 6]$ and $[4, 10]$, the three bounded slabs hold $2 + 3 + 2 = 7$ entries for three segments. Persistent search trees cut storage to $O(n)$ by sharing structure between consecutive slabs, which differ by only a few insertions and deletions, an elegant result that motivates the alternative below.

## Trapezoidal maps

A trapezoidal map shoots vertical extensions from each endpoint only up and down to the nearest segment, instead of across the whole plane. Inside a bounding box, $n$ non-crossing segments produce at most $3n + 1$ trapezoids and at most $6n + 4$ vertices, so the map has linear size. Built by randomized incremental insertion, with a search structure (a directed acyclic graph of x-tests and segment tests) recording the history, it takes $O(n \log n)$ expected time and answers queries in $O(\log n)$ expected time. Example: 100 segments give at most 301 trapezoids. The same decomposition into trapezoids is used by Boolean engines and mask fracturers, so the structure doubles as a geometric representation.

## Going deeper

de Berg et al., "Computational Geometry: Algorithms and Applications" (3rd ed., 2008), chapter 6; Sarnak and Tarjan, "Planar point location using persistent search trees" (Communications of the ACM, 1986); Seidel, "A simple and fast incremental randomized algorithm for computing trapezoidal decompositions" (Computational Geometry, 1991).
