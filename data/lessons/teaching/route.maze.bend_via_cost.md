# Bend and via costs in maze routing

Shortest is not the same as best. In a stack where each layer carries one direction, every bend in a route forces a via, and vias add resistance, block tracks on the layers they cross and reduce yield. Routers therefore price bends and vias explicitly, which changes both the search state and the routes it finds.

## Pricing a route

A route's cost is the sum of its step costs plus a penalty per bend or via: $cost = L + b \cdot B$ for length $L$, bend count $B$ and bend cost $b$, with more terms for wrong-way steps or congested edges. In a preferred-direction stack a bend in the plane is a layer change, so bend cost and via cost are the same thing seen from 2D.

Example: moves N N N N E E E N N N N E E E N have length 15 and change direction 4 times. With a bend cost of 6 the route costs $15 + 24 = 39$. A route with one more step and two fewer bends would cost $16 + 12 = 28$, so the longer route wins. The bend penalty encodes the router's real preference for straight wires.

## Direction-aware search

A bend cost depends on the direction a route arrives from, so a search on cells alone cannot price it. The router searches over states (cell, incoming direction): moving in the same direction costs one step, turning costs one step plus the bend cost. Dijkstra or A* on this expanded graph, four states per cell on a 2D grid, finds the exact cheapest route.

Example: on a grid where the shortest route has 9 steps and 4 bends and another has 11 steps and 2 bends, a bend cost of 4 gives $9 + 16 = 25$ against $11 + 8 = 19$, so the search returns the longer, straighter route at cost 19 (or something cheaper still if it exists). In 3D the state is simply the (cell, layer) vertex, and the via edge carries the cost; the effect is the same.

## When a longer route wins

Two candidates make the trade explicit. Route A is $L_A$ steps with $k_A$ vias, route B is $L_B > L_A$ steps with $k_B < k_A$ vias. With step cost 1 and via cost $v$, B is cheaper when $L_B + k_B v < L_A + k_A v$, that is $v > (L_B - L_A)/(k_A - k_B)$.

Example: A is 11 steps with 5 vias, B is 21 steps with 4 vias. B wins only when $v > 10$, so at any realistic via cost the short route stays. Choosing the via cost is a policy decision: too low and routes zigzag through via stacks that hurt resistance and yield; too high and nets take long detours that add capacitance and congestion. It is exactly the kind of knob that belongs in a router's parameters rather than its core.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.).
