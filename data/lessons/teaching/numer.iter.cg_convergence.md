# CG convergence and the condition number

How many CG iterations a power-grid solve needs is set by the condition number of the matrix. The classical bound converts $\kappa$ into an iteration count, and the way $\kappa$ grows with mesh refinement explains why large grids get slow and why preconditioners are not optional. Tool designers use this estimate to budget run time.

## The square-root-of-kappa bound

For SPD $A$ with condition number $\kappa = \lambda_{max}/\lambda_{min}$, CG satisfies
$$\|e_k\|_A \le 2\left(\frac{\sqrt\kappa - 1}{\sqrt\kappa + 1}\right)^k\|e_0\|_A.$$
With $\kappa = 10^4$, $\sqrt\kappa = 100$ and the factor is $99/101 = 0.9802$ per iteration. Steepest descent has factor $(\kappa - 1)/(\kappa + 1) = 0.9998$, which is why CG is so much faster: it depends on $\sqrt\kappa$, not $\kappa$.

The bound comes from Chebyshev polynomials on the interval $[\lambda_{min}, \lambda_{max}]$. It is often pessimistic, since clustered or outlying eigenvalues help CG, but it gives the right scaling.

## Iterations needed for a tolerance

Solving $2\rho^k \le \varepsilon$ gives
$$k \ge \frac{\ln(\varepsilon/2)}{\ln\rho} \approx \frac{\sqrt\kappa}{2}\ln\frac{2}{\varepsilon}.$$
For $\kappa = 10^4$ and $\varepsilon = 10^{-6}$: $\ln(5 \times 10^{-7}) = -14.51$ and $\ln 0.9802 = -0.0200$, so $k = 726$ iterations. Tightening the tolerance to $10^{-9}$ adds only about $50\ln 1000 = 345$ more: the cost grows with the logarithm of the accuracy.

A simulator's run-time estimate is therefore iterations times SpMV cost. For $n = 10^7$ unknowns, 726 iterations of a 0.1 s SpMV is over a minute, a figure that motivates preconditioning and GPU kernels.

## Why power grids are ill-conditioned

A grid Laplacian with spacing $h$ has eigenvalues from about $\lambda_{min} \sim h^2$ (relative to the largest, for a domain held at its boundary by pads) to $\lambda_{max} \sim 1$, so $\kappa = O(h^{-2})$. For a $k$-by-$k$ grid, $\kappa$ grows like $k^2$, and CG iterations like $\sqrt\kappa \sim k$.

So refining a grid by a factor $r$ in each direction multiplies unknowns by $r^2$, $\kappa$ by $r^2$ and iterations by $r$: total work grows like $r^3$. A solve taking 200 iterations on a $1000 \times 1000$ grid takes about 400 on $2000 \times 2000$, each iteration four times as expensive. Sparse pad placement makes things worse, since nodes far from any pad are weakly held. This is why IR-drop solvers use incomplete Cholesky or algebraic multigrid preconditioners, which aim to make the iteration count nearly independent of grid size.

## Going deeper

Greenbaum, "Iterative Methods for Solving Linear Systems" (SIAM, 1997), chapter 3; Trefethen and Bau, "Numerical Linear Algebra" (SIAM, 1997), lecture 38.
