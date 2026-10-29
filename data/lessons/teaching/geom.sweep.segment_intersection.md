# Line-segment intersection by sweep

Finding all crossings among wire segments is the textbook sweep problem and a direct layout need: crossings on the same layer are shorts, and crossings between layers are where vias may go and where coupling is computed. For Manhattan segments the sweep is especially simple; for general segments the Bentley-Ottmann algorithm finds all $k$ crossings among $n$ segments in $O((n + k) \log n)$ time.

## Crossings of horizontal and vertical segments

A horizontal segment from $(x_1, y)$ to $(x_2, y)$ and a vertical segment from $(x, y_1)$ to $(x, y_2)$ cross when $x_1 < x < x_2$ and $y_1 < y < y_2$, a pair of interval tests. Example: H from $(0, 3)$ to $(8, 3)$ and V from $(5, 0)$ to $(5, 6)$ cross at $(5, 3)$; V from $(9, 0)$ to $(9, 6)$ does not, since 9 is outside $[0, 8]$. Testing every pair costs $O(n_h n_v)$. Whether endpoints touching a segment count as crossings is again a predicate choice; for shorts detection on merged polygons touching is connectivity, while for crossing counts it is usually excluded. Manhattan crossings always lie on existing coordinates, so the computation is exact. On a large layer the pair loop is far too slow, which is what the sweep below fixes.

## The orthogonal sweep

Sweep left to right. Events are the endpoints of horizontal segments and the x-positions of vertical segments. Insert a horizontal segment into a tree keyed by y at its left endpoint and delete it at its right endpoint. At a vertical segment, run a one-dimensional range query on the tree for y in $[y_1, y_2]$, reporting every active horizontal it crosses. Example: at a vertical segment at $x = 5$, three horizontals may be active (their x-ranges contain 5), and a range query over the vertical's span reports the two whose y lies inside it. Each query costs $O(\log n + k_v)$, so the whole algorithm runs in $O(n \log n + k)$, which is optimal for reporting.

## Bentley-Ottmann and output sensitivity

For segments at arbitrary angles, the status is the list of segments ordered by their y at the sweep line, and that order changes exactly at crossings. Bentley-Ottmann keeps an event queue of endpoints and of crossings between status neighbours, discovered on the fly; only adjacent segments can cross next. Each of the $2n + k$ events costs $O(\log n)$, giving $O((n + k)\log n)$. Compare brute force at $n(n-1)/2$ pair tests. With $n = 2^{16}$ and $k = n$, brute force does about $2.1 \times 10^9$ tests while the sweep does about $2.1 \times 10^6$ operations, a thousand times fewer. When $k$ approaches $n^2$ the advantage disappears, but layout data is sparse: most segments cross few others.

## Going deeper

Bentley and Ottmann (1979), as above; de Berg et al., "Computational Geometry: Algorithms and Applications" (3rd ed., 2008), chapter 2; Chazelle and Edelsbrunner, "An optimal algorithm for intersecting line segments in the plane" (Journal of the ACM, 1992).
