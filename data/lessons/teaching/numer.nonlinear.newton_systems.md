# Newton-Raphson for nonlinear systems

A DC operating point, and every implicit timestep of a transient, is a system of nonlinear equations $F(x) = 0$, one KCL equation per node with device currents that depend nonlinearly on node voltages. Newton-Raphson for systems solves it by repeated linearization, and each iteration is one sparse linear solve. It is the outer loop that drives almost all of the linear algebra in a circuit simulator.

## The Jacobian and the linearized system

Expand $F$ around the current iterate: $F(x_k + \Delta) \approx F(x_k) + J(x_k)\Delta$, with Jacobian entries $J_{ij} = \partial F_i/\partial x_j$. Setting the linear model to zero gives the Newton system
$$J(x_k)\,\Delta = -F(x_k), \qquad x_{k+1} = x_k + \Delta.$$
For $F_1 = x^2 + 2y - 5$, $F_2 = x + y^2 - 3$:
$$J = \begin{pmatrix} 2x & 2 \\ 1 & 2y \end{pmatrix}, \qquad \det J = 4xy - 2.$$
At $(1.5, 1.5)$, $\det J = 7$. In a circuit, $J$ is the MNA matrix of the circuit with every nonlinear device replaced by its small-signal conductances at the present operating point, which is why it has the same sparse pattern at every iteration.

## One Newton step on a two-equation system

Continue from $(x_0, y_0) = (1.5, 1.5)$. The residual is $F_1 = 2.25 + 3 - 5 = 0.25$ and $F_2 = 1.5 + 2.25 - 3 = 0.75$. Solve
$$\begin{pmatrix} 3 & 2 \\ 1 & 3 \end{pmatrix}\Delta = \begin{pmatrix} -0.25 \\ -0.75 \end{pmatrix}.$$
By Cramer's rule with $\det = 7$: $\Delta_x = (-0.25 \cdot 3 - 2 \cdot (-0.75))/7 = 0.75/7 = 0.1071$ and $\Delta_y = (3 \cdot (-0.75) - (-0.25) \cdot 1)/7 = -2/7 = -0.2857$. So $(x_1, y_1) = (1.6071, 1.2143)$.

In a simulator the solve is a sparse LU refactorization with the pattern from the symbolic phase, followed by triangular solves. The cost of a Newton iteration is dominated by device evaluation (filling $J$ and $F$) and that factorization.

## Residual norms and when to stop

At $(1.6071, 1.2143)$ the residual is $F_1 = 2.5829 + 2.4286 - 5 = 0.0115$ and $F_2 = 1.6071 + 1.4745 - 3 = 0.0816$, so $\|F\|_\infty$ fell from 0.75 to 0.082. Near a root with nonsingular Jacobian, Newton converges quadratically: the error roughly squares each step, so the number of correct digits doubles.

A solver stops when both the update and the residual are small. The update test alone can be fooled by a stalled iteration, and the residual test alone by badly scaled equations, so SPICE-style simulators require both, with tolerances relative to the size of the unknowns and currents. Far from the root none of this is guaranteed: Newton can diverge or cycle, which is what damping, limiting and continuation are for.

## Going deeper

Kelley, "Solving Nonlinear Equations with Newton's Method" (SIAM, 2003); Najm, "Circuit Simulation" (Wiley, 2010), chapter 3.
