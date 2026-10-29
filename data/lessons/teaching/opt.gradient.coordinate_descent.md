# Coordinate descent

Coordinate descent minimizes over one variable at a time, holding the others fixed. For a quadratic each one-variable minimization is a closed-form update, and a full sweep is exactly one Gauss-Seidel iteration on the linear system $\nabla f = 0$. In placement this is the classic move-each-cell-to-the-weighted-centre-of-its-neighbours iteration, and in detailed placement and sizing it appears as cell-by-cell or gate-by-gate improvement that needs no global solve.

## A single coordinate update

For $f(x, y) = \tfrac{1}{2}(ax^2 + 2hxy + by^2) - ux - vy$, holding y fixed and setting $\partial f / \partial x = ax + hy - u = 0$ gives
$$x \leftarrow \frac{u - hy}{a}.$$
Example: $f = \tfrac{1}{2}(5x^2 - 2xy + 4y^2) + 2x + 3y$, so $a = 5$, $h = -1$, $u = -2$. From $y = 1$: $x = (-2 - (-1)(1))/5 = -0.2$.

Each update is an exact line search along a coordinate axis, so the objective never increases. The update reads only the variables coupled to x, its neighbours in the netlist, which makes the method cheap per step and easy to run incrementally after a small netlist or placement change. The weakness is directional: it can only move along axes, so it is slow when the coupling h is strong and the valley runs diagonally.

## A full sweep is Gauss-Seidel

A sweep updates every coordinate in turn, each using the newest values. For the two-variable quadratic: first $x_1 = (u - hy_0)/a$, then $y_1 = (v - hx_1)/b$. Example: $f = \tfrac{1}{2}(2x^2 - 2xy + 4y^2) + 7x + y$ from $(4, 3)$, so $a = 2$, $b = 4$, $h = -1$, $u = -7$, $v = -1$. Then $x_1 = (-7 + 3)/2 = -2$ and $y_1 = (-1 - 2)/4 = -0.75$.

These are exactly the Gauss-Seidel updates for the system $\begin{pmatrix} a & h \\ h & b \end{pmatrix}\begin{pmatrix} x \\ y \end{pmatrix} = \begin{pmatrix} u \\ v \end{pmatrix}$. Gauss-Seidel converges for every symmetric positive definite matrix, so coordinate descent on a strictly convex quadratic always converges, at a rate set by the conditioning. Using the old values for all coordinates instead gives the Jacobi method, which parallelizes better but may converge more slowly.

## Cell-by-cell placement relaxation

In quadratic placement the coordinate update for a cell is the weighted average of its neighbours' current positions. Take a chain pad(0), c1, c2, c3, pad(79) with net weights 1, 3, 2, 2, starting at c1 = 17, c2 = 38, c3 = 33. One sweep: $c_1 = (1 \cdot 0 + 3 \cdot 38)/4 = 28.5$; $c_2 = (3 \cdot 28.5 + 2 \cdot 33)/5 = 30.3$; $c_3 = (2 \cdot 30.3 + 2 \cdot 79)/4 = 54.65$.

Information travels one cell per sweep along the chain, so a long chain needs many sweeps: that is the poor conditioning of placement Laplacians seen from the inside. Multigrid and conjugate gradient fix it globally. Coordinate-wise moves remain valuable in detailed placement, where each cell or small window is optimized exactly against its fixed neighbours, as in optimal-region and window-based refinement.

## Going deeper

Wright, "Coordinate descent algorithms" (Mathematical Programming, 2015); Saad, "Iterative Methods for Sparse Linear Systems" (2nd ed., 2003), chapter 4 on Jacobi and Gauss-Seidel.
