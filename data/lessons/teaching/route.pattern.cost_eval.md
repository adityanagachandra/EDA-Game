# Choosing patterns on a congestion map

A pattern router does not ask which shape is shortest; all monotone shapes have the same length. It asks which is cheapest on the current congestion map, where each gcell edge's cost rises as its demand approaches capacity. Evaluating patterns against that map, with a penalty for bends, is the inner loop of fast global routing.

## The cheaper L

Place the pins at gcells $(0, 0)$ and $(d_x, d_y)$. Let $h_{y,x}$ be the cost of the horizontal edge from $(x, y)$ to $(x+1, y)$ and $v_{x,y}$ the cost of the vertical edge from $(x, y)$ to $(x, y+1)$. The horizontal-first L costs the row-0 edges plus the column-$d_x$ edges; the vertical-first L costs the column-0 edges plus the row-$d_y$ edges.

Example: pins at (0, 0) and (2, 1); row 0 costs 3 and 9, column 2 costs 8, so HV costs 20; column 0 costs 8 and row 1 costs 1 and 9, so VH costs 18. The router picks VH at 18. Only the boundary rows and columns of the bounding box matter for L shapes, which is what makes them so cheap to evaluate.

## Adding Z shapes

An HVH Z with its jog at column $j$ costs row 0 from 0 to $j$, column $j$ from 0 to $d_y$, and row $d_y$ from $j$ to $d_x$. A VHV Z with its jog at row $i$ costs column 0 up to $i$, row $i$ across, and column $d_x$ from $i$ up. The router evaluates every jog position and keeps the cheapest pattern overall.

With prefix sums along each row and column, each Z costs three lookups, so all $d_x + d_y$ patterns are evaluated in linear time. The Z shapes let a route slip through a cheap column or row in the middle of the box while still avoiding expensive edges at the boundary. On a map where both L shapes cross a hot spot, a Z often wins at the price of one extra bend.

## Bends as a cost

Each bend becomes a via after layer assignment, so pattern routers add a penalty $b$ per bend: an L costs its edges plus $b$, a Z its edges plus $2b$. A cheap Z can then lose to a slightly more expensive L.

Example: the best L costs 7 in edges and the best Z costs 18; with $b = 2$ the L costs 9 and the Z 22. With edge costs closer together, say L at 14 and Z at 11, a penalty of $b = 4$ flips the choice to the L (18 against 19). The bend penalty is a policy knob: it trades congestion relief against via count, and its right value depends on the stack and on how expensive vias are in yield and resistance.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.); the FastRoute papers by Pan and Chu.
