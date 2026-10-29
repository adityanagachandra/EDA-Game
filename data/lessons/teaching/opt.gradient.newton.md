# Newton's method for minimization

Newton's method for minimization replaces the function by its second-order Taylor model and jumps to that model's minimum. Near a nondegenerate minimum it converges quadratically, doubling the correct digits each step, and on a quadratic it lands in one step. Its price is the Hessian, too expensive for full-chip placement but practical for small dense subproblems such as sizing a gate cluster, solving a nonlinear circuit DC point or refining a few macro positions.

## The Newton step in one dimension

The quadratic model of f at $x_k$ is $f(x_k) + f'(x_k)s + \tfrac{1}{2}f''(x_k)s^2$, minimized at $s = -f'(x_k)/f''(x_k)$ when $f''(x_k) > 0$. So
$$x_{k+1} = x_k - \frac{f'(x_k)}{f''(x_k)},$$
the root-finding Newton iteration applied to $f'$.

Example: $f(x) = e^x - 8x$, so $f'(x) = e^x - 8$ and $f''(x) = e^x$, and $x_{k+1} = x_k - 1 + 8e^{-x_k}$. From $x_0 = 0.5$: $x_1 = 0.5 - 1 + 8 \cdot 0.6065 = 4.352$, overshooting the true minimizer $\ln 8 = 2.079$, because far from the solution the quadratic model is poor. The next steps are 3.455, 2.708, 2.241, 2.092 and 2.080: once close, the error roughly squares at each step. Damped Newton, which takes a fraction of the step chosen by line search, is how solvers stay safe while far away.

## One step on a quadratic

In several variables the step solves $H s = -\nabla f$. For a quadratic the Taylor model is exact, so one step from any start reaches the minimizer, which is the solution of $\nabla f = 0$.

Example: $f(x, y) = 2x^2 - 2xy + 2y^2 + 10x + 12y$. The gradient is $(4x - 2y + 10,\ -2x + 4y + 12)$ and the Hessian is $\begin{pmatrix} 4 & -2 \\ -2 & 4 \end{pmatrix}$, positive definite. Setting the gradient to zero: $4x - 2y = -10$ and $-2x + 4y = -12$, giving $(x, y) = (-5.333, -5.667)$, whatever the start. This is why quadratic placement is solved by a linear system rather than by iterating descent, and why Newton-type methods dominate when the function is nearly quadratic near its optimum. The cost is forming and factoring H: $O(n^3)$ for dense matrices.

## Negative curvature and safeguards

Newton's method goes to a stationary point of the model, which is a minimum only when the Hessian is positive definite. With negative curvature it heads for a maximum or a saddle. Take $f(x) = x^3 - 3x$, with $f'(x) = 3x^2 - 3$ and $f''(x) = 6x$; the step is $x_{k+1} = x_k - (x_k^2 - 1)/(2x_k)$. From $x_0 = -0.5$: $f'' = -3 < 0$, and $x_1 = -0.5 - (-0.75)/(-1) = -1.25$, moving toward $x = -1$, the local maximum.

Practical codes guard against this. A modified Newton method adds a multiple of the identity to H until it is positive definite. Trust-region methods minimize the model only within a radius where it is trusted. Line-search methods check that the step is a descent direction, $\nabla f^T s < 0$, and fall back to the negative gradient when it is not. Circuit simulators apply the same idea to Newton-Raphson DC solves with limiting and damping.

## Going deeper

Nocedal and Wright, "Numerical Optimization" (2nd ed., 2006), chapters 3 and 4 on Newton and trust-region methods; Boyd and Vandenberghe, "Convex Optimization" (2004), section 9.5.
