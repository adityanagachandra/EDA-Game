# Force-directed placement

Force-directed placement models connections as springs that pull connected cells together and, in its modern forms, density as a force that pushes cells out of crowded regions. It is the physical picture behind quadratic placement and the ancestor of electrostatic placement. The classic algorithm moves one cell at a time to the point where its spring forces balance.

## Springs and the zero-force target

Treat each connection of weight $c_j$ from cell $i$ to cell $j$ as a spring with force $c_j$ times the distance, Hooke's law. The total force on cell $i$ at position $(x_i, y_i)$ is $F_x = \sum_j c_j(x_j - x_i)$ and similarly for y. It vanishes at the weighted centroid of the neighbours:
$$x_i^0 = \frac{\sum_j c_j x_j}{\sum_j c_j}, \qquad y_i^0 = \frac{\sum_j c_j y_j}{\sum_j c_j}.$$
For neighbours at (26, 20) with weight 2 and (10, 39) with weight 1, the zero-force target is (20.67, 26.33). This is exactly the minimiser of the quadratic wirelength for that one cell, which is why force-directed and quadratic placement are two views of the same objective.

## Net force on a cell

At any position the net force tells which way and how hard a cell is pulled. For a cell at (25, 19) with neighbours at x = 17, 26 and 40 and weights 3, 4 and 4, $F_x = 3(17 - 25) + 4(26 - 25) + 4(40 - 25) = -24 + 4 + 60 = 40$, a strong pull to the right. Force magnitude is a natural priority: the iterative algorithm picks the cell with the largest force, or the most connections, and moves it first. Forces also explain why quadratic placement clusters: the springs pull everything toward the centre and nothing pushes back, so modern force-directed placers such as Kraftwerk add a density force computed from the placement's density distribution, which pushes cells from dense to sparse regions.

## The ripple-move algorithm

The classic iterative force-directed algorithm proceeds cell by cell. Pick a cell, compute its zero-force target and round it to the nearest legal grid slot. If the slot is free, move there and lock the cell. If it holds another cell, move in anyway and pick up the displaced cell, which then looks for its own zero-force target, a ripple of moves. If the target is the cell's current location, it stays. Abort conditions stop endless rippling, for example after too many moves without improvement. The rounding step makes this a discrete local search: a cell at (35, 32) with a target that rounds to (17, 38) moves a Manhattan distance of 24. The method is simple and intuitive but slow and greedy, which is why large-scale placers moved to simultaneous analytical methods.

## Going deeper

Quinn and Breuer, "A force directed component placement procedure for printed circuit boards" (IEEE Trans. Circuits and Systems, 1979); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), section 4.3.3.
