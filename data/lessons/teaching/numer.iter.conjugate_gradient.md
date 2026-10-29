# The conjugate gradient method

Conjugate gradients (CG) is the iterative method for symmetric positive definite systems, the matrices of resistive power grids, RC interconnect and thermal networks. It needs only matrix-vector products and a few vectors, converges far faster than relaxation, and in preconditioned form is the core of most large IR-drop and thermal solvers.

## Minimizing the energy along conjugate directions

For SPD $A$, solving $Ax = b$ is equivalent to minimizing the energy
$$\phi(x) = \tfrac12 x^TAx - b^Tx,$$
whose gradient is $Ax - b = -r$. For a power grid, $\phi$ is half the dissipated power minus the work done by the injected currents. Example: $A = \begin{pmatrix} 4 & 1 \\ 1 & 3 \end{pmatrix}$, $b = (1, 2)$, $x = (1, 0)$: $\phi = \frac12 \cdot 4 - 1 = 1$.

Steepest descent moves along $r$ and zig-zags when $A$ is ill-conditioned. CG instead uses search directions that are $A$-conjugate, $p_i^TAp_j = 0$, so minimizing along a new direction never undoes progress along earlier ones. After $k$ steps, $x_k$ minimizes $\phi$, equivalently the $A$-norm of the error, over the Krylov space $x_0 + \text{span}\{r_0, Ar_0, \dots, A^{k-1}r_0\}$.

## One CG iteration by hand

From $x_0 = 0$: $r_0 = b$, $p_0 = r_0$. Then
$$\alpha_0 = \frac{r_0^Tr_0}{p_0^TAp_0}, \quad x_1 = x_0 + \alpha_0p_0, \quad r_1 = r_0 - \alpha_0Ap_0,$$
$$\beta_0 = \frac{r_1^Tr_1}{r_0^Tr_0}, \quad p_1 = r_1 + \beta_0p_0.$$
For $A = \begin{pmatrix} 4 & 1 \\ 1 & 3 \end{pmatrix}$ and $b = (1, 2)$: $r_0^Tr_0 = 5$, $Ap_0 = (6, 7)$, $p_0^TAp_0 = 20$, so $\alpha_0 = 0.25$, $x_1 = (0.25, 0.5)$ and $r_1 = (1 - 1.5, 2 - 1.75) = (-0.5, 0.25)$. Check: $r_1^Tr_0 = -0.5 + 0.5 = 0$, the residuals are orthogonal.

Each iteration costs one SpMV, two dot products and three vector updates, and stores four vectors. The dot products are global reductions, which matter on parallel machines.

## Exact termination and the Krylov subspace

In exact arithmetic CG finds the exact solution in at most $n$ iterations, and in fact in at most as many iterations as $A$ has distinct eigenvalues: the error after $k$ steps is $q(A)e_0$ for the best polynomial $q$ of degree $k$ with $q(0) = 1$, and a polynomial with roots at all distinct eigenvalues annihilates the error. A matrix with eigenvalues only 2 and 5 is solved in 2 iterations whatever its size.

In floating point, conjugacy is lost gradually and termination is not exact, so CG is used as an iterative method stopped by a residual test. The polynomial view explains the practical behaviour: convergence is fast when eigenvalues are clustered, which is precisely what a good preconditioner arranges.

## Going deeper

Shewchuk, "An Introduction to the Conjugate Gradient Method Without the Agonizing Pain" (CMU technical report, 1994); Saad, "Iterative Methods for Sparse Linear Systems" (2nd ed., SIAM, 2003), chapter 6.
