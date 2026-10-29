# Lee's maze router

Lee's algorithm (1961) is the founding maze router: breadth-first search on the routing grid, phrased as a wavefront that floods outward from the source until it reaches the target. It is guaranteed to find a shortest route if one exists. Every later maze router, from Hadlock's to A*, is a way to get Lee's guarantee while expanding fewer cells or using less memory.

## Wavefront expansion

Label the source 0. Every free neighbour of a cell labelled $k$ that has no label yet receives $k + 1$. The cells labelled $k$ form the $k$-th wavefront, and they are exactly the cells at shortest-route distance $k$ from the source, because blockages are simply never labelled and the wave flows around them. The search stops when the target receives a label $D$, the length of the shortest route in grid steps.

If the Manhattan distance between the pins is 6 but a wall forces the route around one end, the target might receive label 8. The cost of the method is the expansion: with unit costs Lee processes each reachable cell once, $O(N)$ for $N$ grid cells. It is complete (it finds a route whenever one exists) and optimal for unit edge costs.

## How many cells the wave touches

Lee stops as soon as the target is labelled, but by then every cell with a smaller label has been labelled too. So the work is at least the number of reachable cells with label below $D$. On an open grid that number grows like the area of a diamond of radius $D$, about $2D^2$ cells, while the route itself has only $D$ steps. Most of the labelled cells lie behind the source or off to the side, nowhere near the target.

Example: on a 6 by 6 grid with a wall, a target at label 8 might be reached only after 20 cells with labels 0 to 7 have been labelled. This waste is what Hadlock's detour ordering and A*'s heuristic attack, and it is why practical routers bound the search to a box around the pins.

## Label memory

A naive implementation stores the full label in every cell. A label can reach the number of cells minus 1, and two more codes are needed for empty and blocked, so a grid of $N$ cells needs $\lceil \log_2 (N + 2) \rceil$ bits per cell: 23 bits for a 2100 by 2400 grid, about 116 million bits in all.

The backtrace only needs to recognise a cell's predecessor, whose label is one less. Storing labels modulo 3 (sequence 1, 2, 3, 1, 2, 3) suffices, because a cell's neighbours on shortest routes carry labels $k-1$, $k$ or $k+1$, which are distinct mod 3. With empty and blocked that is 5 codes, 3 bits. Akers' coding (sequence 0, 0, 1, 1, 0, 0, 1, 1) needs only two label values, so 4 codes fit in 2 bits per cell. Memory dropped by an order of magnitude with no loss of optimality.

## Going deeper

C. Y. Lee, "An Algorithm for Path Connections and Its Applications" (1961); S. B. Akers, "A Modification of Lee's Path Connection Algorithm" (1967); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.).
