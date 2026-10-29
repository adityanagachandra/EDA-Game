# The plane-sweep paradigm

A plane sweep turns a two-dimensional problem into a sequence of one-dimensional ones. A vertical line moves across the plane, stopping only at events (shape edges, endpoints, intersections), while a status structure holds the objects the line currently cuts. Union area, overlap reporting, segment intersection, Boolean operations and connectivity extraction in layout tools are all sweeps. The paradigm is why these problems cost $O(n \log n)$ plus output, not $O(n^2)$.

## Events and stops

For rectangles, the events are their left and right edges: $2n$ events for $n$ rectangles, sorted by x in $O(n \log n)$. Nothing changes between consecutive event x-values, so the sweep only needs to stop at distinct event coordinates, processing all events at one x together. Example: rectangles $[0, 5] \times \cdot$, $[3, 9] \times \cdot$ and $[5, 12] \times \cdot$ have six events at x-values $0, 3, 5, 5, 9, 12$, so five stops. The order of events at the same x is a design decision with consequences: processing removals before insertions makes touching shapes non-overlapping, and the reverse makes them overlap. Engines choose the order to match the predicate they implement.

## The status structure

Between stops, the status holds the objects the sweep line cuts, ordered along the line, typically in a balanced search tree keyed by y. At a left edge the rectangle is inserted; at a right edge it is deleted; each update costs $O(\log m)$ for $m$ active objects. The maximum status size is the depth of the deepest overlap of x-ranges, which for layout data is usually far smaller than $n$: a vertical line across a chip cuts only the shapes in one column. Example: x-ranges $[0, 6]$, $[2, 8]$, $[4, 10]$, $[7, 12]$ are three deep at $x = 5$ and never four deep, so the status never holds more than three rectangles. Memory for a sweep is proportional to that peak, not to the whole layout.

## Why the sweep is fast

A sweep only compares objects that are in the status at the same time, which means their x-ranges overlap. Pairs whose x-ranges are disjoint are never compared at all. Brute force tests all $\binom{n}{2}$ pairs; the sweep tests only x-overlapping pairs and, with a good status structure, only those that also come close in y. Example: six rectangles have 15 pairs, but if only 5 of them overlap in x, the sweep hands just those 5 to the y-test. The total cost becomes $O(n \log n)$ for sorting and status updates plus $O(k)$ or $O(k \log n)$ for the $k$ reported results, an output-sensitive bound. This locality is what makes full-chip checks tractable.

## Going deeper

Bentley and Ottmann, "Algorithms for reporting and counting geometric intersections" (IEEE Transactions on Computers, 1979); de Berg et al., "Computational Geometry: Algorithms and Applications" (3rd ed., 2008), chapter 2.
