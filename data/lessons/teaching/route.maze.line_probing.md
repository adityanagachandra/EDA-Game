# Line-probe routing

Line-probe routers, by Mikami and Tabuchi (1968) and Hightower (1969), search with straight lines instead of cells. They shoot escape lines from the pins, then perpendicular lines from points on those, until a line from the source meets a line from the target. They store a few segments instead of a label per cell, and the routes they find have few bends.

## Minimum bends

Mikami-Tabuchi starts with level-0 lines: from each pin, straight lines in all four directions until they hit a blockage or the grid edge. At level $i + 1$ it launches perpendicular lines from every grid point of the level-$i$ lines. The search stops when a source line intersects a target line. Every level adds one bend, so the first intersection gives a route with the fewest possible bends, which is also the fewest straight segments.

Example: if the target sits behind a wall with a gap, the route might need 3 bends: out of the source, into the gap's row, through, and down to the target. The route is not necessarily shortest in length, but in a preferred-direction stack fewer bends mean fewer vias.

## Probe lines and what they cover

A level-0 line runs from the pin until the first blocked cell or the grid edge. The cells on the four level-0 lines of the source are exactly the cells reachable with zero bends. The cells reachable with at most one bend are those on the level-0 lines plus every cell on a perpendicular line launched from them.

Example: a source in an open corner might cover 8 cells with its level-0 lines and 15 with one bend. A line-probe router stores each line as a segment (two endpoints), so memory scales with the number of lines rather than with the area of the grid. Hightower's variant is greedier: it launches only a few lines from carefully chosen escape points, which is faster but can fail to find a route that exists.

## Line probing versus maze routing

Maze routers label cells and guarantee shortest routes; line probers store segments and aim for few bends. Mikami-Tabuchi is complete (it finds a route if one exists) because every grid point of every line launches the next level. Hightower trades completeness for speed. Neither handles general edge costs as naturally as a priority search.

In practice, line-probe ideas survive in gridless detailed routers and in fast connection searches through open areas, often followed by a maze search when the probes fail. Because each line is a segment, the representation pairs well with interval and corner-stitched data structures. The fewest-bend property also matches what a preferred-direction router wants: long straight runs with few layer changes.

## Going deeper

K. Mikami and K. Tabuchi, "A Computer Program for Optimal Routing of Printed Circuit Conductors" (1968); D. W. Hightower, "A Solution to Line-Routing Problems on the Continuous Plane" (1969); Sherwani, "Algorithms for VLSI Physical Design Automation".
