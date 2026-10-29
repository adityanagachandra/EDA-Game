# Linear programming and its geometry

A linear program minimizes or maximizes a linear objective over a region cut out by linear inequalities. LPs sit inside EDA tools everywhere: linearized sizing and budgeting, clock-skew scheduling, LP relaxations of routing and covering ILPs, and the bounding step of every branch-and-bound solver. The geometry of two-variable LPs explains the whole theory: the feasible region is a convex polygon and an optimum, when one exists, sits at a corner.

## Vertices from tight constraints

Each inequality $a_i x + b_i y \le r_i$ is a half-plane, and the feasible region is their intersection with $x, y \ge 0$: a convex polygon. A vertex is a point where two boundary lines meet and every other constraint still holds. To find the vertex where constraints (1) and (2) are both tight, solve them as equations.

Take maximize $z = 6x + 5y$ subject to (1) $5x + 2y \le 22$ and (2) $x + 5y \le 9$. Solving $5x + 2y = 22$ and $x + 5y = 9$ gives $x = 4$, $y = 1$ (by elimination or Cramer's rule). The objective there is $6 \cdot 4 + 5 \cdot 1 = 29$. The other vertices come from the axes: $(0, 0)$, the x-intercept of the binding constraint $(4.4, 0)$, and the y-intercept $(0, 1.8)$.

## The optimum is at a vertex

A linear objective has parallel level lines $c_1 x + c_2 y = z$. Sliding the line in the improving direction, the last feasible point it touches is a vertex (or a whole edge, when the objective is parallel to it). So a bounded, feasible LP has an optimal vertex, and checking every vertex solves a small LP.

For the example above the vertices are $(0, 0)$ with $z = 0$, $(4.4, 0)$ with $z = 26.4$, $(4, 1)$ with $z = 29$, and $(0, 1.8)$ with $z = 9$. The optimum is 29 at $(4, 1)$. Changing the objective to $z = x + 10y$ moves it to $(0, 1.8)$ with $z = 18$: the region is fixed, the objective picks the corner. The simplex method walks from vertex to adjacent vertex, improving each step, and stops when no neighbour is better. Interior-point methods instead cross the inside of the region; both reach the same optimal value.

## Counting corners, redundancy and degeneracy

Adding a constraint can add a corner, remove corners, or change nothing. Add (3) $y \le 5$ to maximize subject to $3x + 2y \le 18$, $x + 3y \le 20$: the line $y = 5$ slices off the top, and the corners become $(0, 0)$, $(6, 0)$, $(8/3, 5)$ and $(0, 5)$, four in all. A constraint whose line never touches the region is redundant and adds no corner.

To count corners mechanically, intersect every pair of boundary lines (including the axes), keep the points that satisfy all constraints, and merge duplicates. Duplicates arise when three or more lines pass through the same point: degeneracy. Degenerate vertices are why simplex can pivot without moving, and why practical codes need anti-cycling rules. In EDA LPs with thousands of similar constraints, such as uniform track capacities, degeneracy is the norm, not the exception.

## Going deeper

Bertsimas and Tsitsiklis, "Introduction to Linear Optimization" (1997), chapters 1 to 3 on geometry and the simplex method; Chvatal, "Linear Programming" (1983), for degeneracy and cycling.
