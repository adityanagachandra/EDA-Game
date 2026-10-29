# Finite-difference solution of Laplace's equation

The simplest field solver discretizes Laplace's equation on a grid and iterates until every node equals the average of its neighbours. Finite differences are rarely the production method for 3D extraction, but they teach the whole pipeline: discretize, solve a sparse linear system, and recover a physical quantity such as capacitance from the field. They are also a good reference implementation for testing faster methods.

## The five-point stencil

In a charge-free, uniform region the potential satisfies $\nabla^2\phi = 0$. On a square grid with spacing $h$, central differences give
$$\frac{\phi_N + \phi_S + \phi_E + \phi_W - 4\phi_0}{h^2} = 0,$$
so every interior node equals the average of its four neighbours. Conductors become nodes held at fixed potentials (Dirichlet boundaries).

Example: neighbours at 1.0, 0.6, 0.8 and 0.2 V give a Jacobi update of $2.6/4 = 0.65$ V. The truncation error of the stencil is $O(h^2)$, so halving the grid spacing cuts the error about fourfold, at four times the unknowns in 2D. Between dielectrics of different permittivity the stencil weights change so that the normal component of $D$ is continuous across the interface.

## Iterating to convergence

The stencil equations form a large sparse linear system. Iterative methods sweep the grid repeatedly. Jacobi updates every node from the previous sweep's values. Gauss-Seidel uses each new value as soon as it is computed, which typically converges about twice as fast; successive over-relaxation (SOR) extrapolates further.

Example on a 1D line, $\phi'' = 0$, boundaries 1 V and 0 V and three interior nodes starting at 0, 0, 0. One Gauss-Seidel sweep left to right gives $x_1 = (1 + 0)/2 = 0.5$, then $x_2 = (0.5 + 0)/2 = 0.25$, then $x_3 = (0.25 + 0)/2 = 0.125$. Jacobi from the same start gives $(0.5, 0, 0)$. Both converge to the straight line $(0.75, 0.5, 0.25)$. Plain Jacobi and Gauss-Seidel slow down on fine grids because smooth errors decay slowly, which is why production solvers use multigrid or preconditioned conjugate gradients.

## From potential to capacitance

Capacitance comes from charge, and charge from Gauss's law: integrate $\varepsilon\,\partial\phi/\partial n$ around a closed contour enclosing the conductor. On a 2D grid of square cells, each cell edge crossed by the contour carries flux $\varepsilon(\Delta\phi/h)h = \varepsilon\Delta\phi$ per unit length, where $\Delta\phi$ is the potential drop across that edge. So
$$C' = \frac{\varepsilon_0\varepsilon_r\sum\Delta\phi}{V}.$$

Example: an inner conductor at 1 V in SiO2 ($\varepsilon_r = 3.9$), with a contour crossing four edges with drops 0.30, 0.25, 0.25 and 0.20 V. Then $Q' = 8.854 \cdot 3.9 \cdot 1.00 = 34.5$ pC/m and $C' = 34.5$ pF/m. Choosing a contour away from sharp corners improves accuracy, since the field is singular at corners. Checking that different contours give the same charge is a simple convergence test.

## Going deeper

Sadiku, "Numerical Techniques in Electromagnetics" (3rd ed.), chapter 3; Saad, "Iterative Methods for Sparse Linear Systems" (2nd ed.), chapter 4.
