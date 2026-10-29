# Manhattan distance and rectilinear geometry

Layout geometry is rectilinear: almost every wire, cell outline and via is a polygon whose edges run horizontally or vertically, on an integer grid. In that world the natural distance between two points is not the straight line but the Manhattan, or $L_1$, distance: the length of the shortest path that moves only along the axes. Wirelength estimates, routing costs, Steiner trees and placement objectives are all built on it, so a tool designer reaches for it before any other metric.

## The L1 metric

For points $p = (x_1, y_1)$ and $q = (x_2, y_2)$ the Manhattan distance is

$$d_1(p, q) = |x_1 - x_2| + |y_1 - y_2|.$$

It is a true metric: zero only for equal points, symmetric, and it obeys the triangle inequality. It is also the exact length of any monotone staircase wire from $p$ to $q$, which is why a router's lower bound for a two-pin net is $d_1$. Example: pins at $(120, 40)$ and $(35, 210)$ nm are $85 + 170 = 255$ nm apart. Its relatives matter too. The Chebyshev distance $d_\infty = \max(|\Delta x|, |\Delta y|)$ appears in square-metric spacing checks, and rotating coordinates by 45 degrees, $u = x + y$, $v = x - y$, turns $d_1$ into $d_\infty$, a trick used by nearest-neighbour and Voronoi algorithms for the $L_1$ metric.

## Manhattan versus Euclidean length

The Euclidean distance $d_2 = \sqrt{\Delta x^2 + \Delta y^2}$ is never larger than $d_1$, and the ratio $d_1/d_2$ ranges from 1 (a purely horizontal or vertical pair) to $\sqrt{2}$ (a 45-degree pair). For pins at $(10, 10)$ and $(40, 50)$ um, $d_1 = 30 + 40 = 70$ um while $d_2 = 50$ um, so the rectilinear wire is 20 um longer than the crow flies. This gap is why rectilinear Steiner trees, not Euclidean ones, are the right model for routed wirelength, and why some advanced nodes allow 45-degree routing on selected layers. A tool should never mix the two metrics silently: a placer that optimizes Euclidean distance while the router pays Manhattan distance is optimizing the wrong thing on diagonal nets.

## Many shortest paths

Between two points there is usually not one shortest rectilinear path but many. On a unit grid with no obstacles, a shortest path from $p$ to $q$ makes exactly $|\Delta x|$ horizontal and $|\Delta y|$ vertical steps, in any order, so the number of shortest paths is the binomial coefficient

$$\binom{|\Delta x| + |\Delta y|}{|\Delta x|}.$$

For $\Delta x = 4$ and $\Delta y = 3$ that is $\binom{7}{4} = 35$ routes of identical length. This freedom is what routers exploit: among equal-length paths they choose by congestion, via count or layer preference. It also explains why a maze router's search can explode on open areas and why routers prefer L- and Z-shaped patterns first: there are only two L-shapes, while the full set of shortest paths grows combinatorially.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapters on global and detailed routing; de Berg, Cheong, van Kreveld and Overmars, "Computational Geometry: Algorithms and Applications" (3rd ed., 2008), chapter 1.
