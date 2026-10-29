# Pattern routing

Maze search is too slow to run on every net of a million-net design. Pattern routing instead tries a handful of predefined shapes for each two-pin connection: the two L shapes with one bend, and the Z shapes with two. Most nets route well this way, and only the leftovers need maze search. Global routers such as FastRoute start every net with patterns.

## Counting Z shapes

Take a two-pin connection whose pins are $d_x \ge 2$ columns and $d_y \ge 2$ rows apart on the gcell grid. A Z shape is a monotone route with exactly two bends. An HVH Z runs horizontally, jogs vertically at some column strictly between the pins, then continues horizontally; there are $d_x - 1$ such columns. A VHV Z jogs horizontally at a row strictly between the pins; there are $d_y - 1$.

So the pattern router considers $(d_x - 1) + (d_y - 1) = d_x + d_y - 2$ Z shapes. Example: pins at (16, 13) and (22, 16) give $5 + 2 = 7$ Z shapes. A jog at a pin's own column would just be an L shape, which is why the jog must be strictly inside.

## All pattern candidates

An L shape has one bend: horizontal then vertical, or vertical then horizontal. When the pins differ in both coordinates there are exactly 2. Together with the Z shapes a pattern router evaluates $2 + (d_x + d_y - 2) = d_x + d_y$ candidates per connection, linear in the distance between the pins.

Each candidate's cost is the sum of the costs of the gcell edges it crosses, plus a bend or via penalty. Computing a cost takes time proportional to the route length, and prefix sums along rows and columns cut that to constant time per candidate. A net in a single row or column has only its straight route. The router picks the cheapest candidate, adds its demand to the edges, and moves on.

## What patterns miss

A monotone route never moves away from the target, and there are $\binom{d_x + d_y}{d_x}$ of them, the lattice-path count. Patterns cover only $d_x + d_y$ of these. Example: for $d_x = 2$, $d_y = 6$ there are $\binom{8}{2} = 28$ monotone routes, of which 2 are L shapes and 6 are Z shapes, so 20 staircase routes with three or more bends are never considered.

That is the trade: patterns are fast and produce few vias, but in a congested region the best route may be a staircase or may need a detour outside the bounding box. Global routers therefore escalate: L and Z first, then monotonic routing by dynamic programming over all staircases, then full maze routing with detours for what is still overflowed. FastRoute later added three-bend patterns as an intermediate step.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.); M. Pan and C. Chu, "FastRoute: A Step to Integrate Global Routing into Placement" (2006).
