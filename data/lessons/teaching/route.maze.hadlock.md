# Hadlock's minimum-detour router

Hadlock's router (1977) keeps Lee's guarantee of a shortest route but changes the order of search. Instead of labelling cells by distance from the source, it labels them by detour number: how many moves away from the target a route needed to reach them. Search proceeds in order of detours, so it pushes straight at the target and spreads sideways only when blocked.

## Counting detours

A move is a detour if it increases the Manhattan distance to the target. Each detour must later be undone by an extra move toward the target, so a route of length $l$ between pins at Manhattan distance $MD$ with $d$ detours satisfies $l = MD + 2d$. The detour count of a route is therefore $d = (l - MD)/2$, however the detours are interleaved.

Example: from (1, 2) to (5, 6) the Manhattan distance is 8. A route of 12 moves, say E E W E W E N E N N E N, makes $(12 - 8)/2 = 2$ detours: the two W moves. Minimising length and minimising detours are the same objective, which is why Hadlock's router is still optimal.

## Expanding by detour number

Hadlock labels each cell with the fewest detours on any route from the source to it, and expands cells in order of that label. A move toward the target keeps the label; a move away increases it by 1. So the search first floods every cell reachable with zero detours, then with one, and so on, until the target is reached. The label of the target is the minimum detour number $d^*$, and the route length is $MD + 2d^*$.

Example: the shortest route around a wall is 11 steps long between pins 3 apart. The minimum detour number is $(11 - 3)/2 = 4$. A cell's Hadlock label is $(dist(s, v) - MD(s, t) + MD(v, t))/2$, where $dist(s, v)$ is its shortest-route distance from the source. This shows the labels are just Lee distances corrected by progress toward the target.

## Why it expands less

With no obstacles, every cell inside the pins' bounding box is reachable with zero detours, and nothing outside it is, so Hadlock expands only the box while Lee floods a diamond of radius $MD$ around the source. Obstacles add detours only where needed. The labels are small integers that grow slowly, which keeps the priority structure simple (a list per detour level).

The price is that the target's position must be known up front, so the method suits two-pin connections; multipin nets are decomposed first. Hadlock is an early instance of goal-directed search: in A* terms, it is Dijkstra on the reduced cost $l - (MD(s,t) - MD(v,t))$, which is exactly twice the detour count.

## Going deeper

F. O. Hadlock, "A Shortest Path Algorithm for Grid Graphs" (1977); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.).
