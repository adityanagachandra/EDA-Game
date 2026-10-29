# Backtrace and path choice in Lee routing

Wavefront labels tell you how far each cell is from the source, but not which route to take. The backtrace recovers a route by walking from the target down the labels. Usually many shortest routes exist, and the choice among them decides bends, vias and congestion. A good backtrace policy is a cheap way to get a better route from the same search.

## Walking down the labels

Start at the target with label $D$. Step to any neighbour labelled $D - 1$, then to a neighbour of that cell labelled $D - 2$, and so on until the source (label 0). Each step lowers the label by exactly one, so the route has exactly $D$ steps and is shortest. The backtrace never needs to search: at every cell at least one neighbour carries the next lower label, because that is how the cell received its own.

The number of choices at the first step is the number of target neighbours labelled $D - 1$. If it is more than one, the router needs a tie-breaking rule. Once the route is found, its cells are marked as used (blocked for the next net) and all other labels are cleared.

## Counting shortest routes

The set of shortest routes forms a directed acyclic graph: an edge from $u$ to $v$ whenever $v$ is a neighbour of $u$ with label one higher. The number of shortest routes to a cell is the sum of the counts of its neighbours with label one less, starting from 1 at the source. Process cells wavefront by wavefront and the count at the target is the number of distinct shortest routes.

On an open grid between pins $d_x$ columns and $d_y$ rows apart this gives the lattice-path count $\binom{d_x + d_y}{d_x}$. Blockages cut it down sharply: behind a wall with a single gap, every shortest route squeezes through the gap, and the count becomes a product of the counts on each side. A large count means freedom, which the backtrace can spend on fewer bends or on avoiding congested cells.

## Tie-breaking for fewer bends

Among shortest routes, a router prefers the one with the fewest bends, because in a preferred-direction stack every bend needs a via. A simple policy keeps the current direction whenever a neighbour in that direction has the next lower label. That greedy rule is not always optimal; the exact minimum comes from dynamic programming over (cell, incoming direction) states on the shortest-route DAG, adding 1 whenever the direction changes.

Example: if the target is reached through a gap in a wall, the best shortest route may need 2 bends (go to the gap's row, cross, turn toward the target), while a careless backtrace zigzags with 5. Tie-breaking costs nothing in search effort, which is why routers add direction and congestion preferences to the backtrace before they reach for heavier tools.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.); Sherwani, "Algorithms for VLSI Physical Design Automation", on maze routing.
