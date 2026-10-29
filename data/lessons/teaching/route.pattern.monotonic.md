# Monotonic routing by dynamic programming

Monotonic routing searches every route that never moves away from the target, not just the L and Z shapes. Because such routes form a directed acyclic grid, dynamic programming finds the cheapest one in time linear in the area of the bounding box. FastRoute 2.0 introduced it as the step between pattern routing and full maze routing.

## The dynamic program

Place the source at $(0, 0)$ and the target at $(d_x, d_y)$. Let $B(x, y)$ be the cheapest monotone cost from the source to gcell $(x, y)$. Then $B(0, 0) = 0$ and
$$B(x, y) = \min\big(B(x-1, y) + h_{y, x-1},\ B(x, y-1) + v_{x, y-1}\big),$$
taking only the terms that exist on the box boundary. Fill the table row by row; $B(d_x, d_y)$ is the answer.

Example: in a 3 by 3 box of gcells the table has 9 entries and each takes one comparison. The route itself is recovered by remembering which term won at each entry. This is the same principle of optimality that underlies shortest paths on a DAG: every monotone route to $(x, y)$ arrives from the west or from the south.

## Scaling to larger boxes

The work is one constant-time update per gcell in the bounding box, $(d_x + 1)(d_y + 1)$ in all, against $\binom{d_x + d_y}{d_x}$ routes enumerated explicitly. For a 4 by 3 box that is 20 table entries against 35 routes; for a 20 by 20 box it is 441 entries against more than $10^{11}$ routes.

Pins can be in any relative position: the router reflects the coordinates so the target is to the north-east, or runs the recurrence in the matching directions. Edge costs come from the current congestion map, so the program is re-run whenever a net is ripped up. Because it never leaves the bounding box, monotonic routing cannot detour around a fully blocked box; that case needs maze routing.

## Monotone versus patterns

Every L and Z shape is monotone, so the best monotone route costs at most the best pattern. The difference is what the extra flexibility buys: a staircase can thread around several expensive edges that every pattern must cross. Example: on a 4 by 3 map, the best pattern might cost 30 while the best monotone route costs 23, a saving of 7.

The cost is more bends. A staircase can turn at every gcell, and every turn is a potential via, so routers either add a bend penalty to the recurrence (which needs the incoming direction in the state, doubling the table) or accept monotone routes only where patterns leave overflow. In the escalation from patterns to monotone to maze, each stage is applied only to the nets the previous stage could not route cleanly.

## Going deeper

M. Pan and C. Chu, "FastRoute 2.0: A High-quality and Efficient Global Router" (2007); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.).
