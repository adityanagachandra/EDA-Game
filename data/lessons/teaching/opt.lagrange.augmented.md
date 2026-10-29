# Augmented Lagrangian method

The augmented Lagrangian method, also called the method of multipliers, adds both a multiplier term and a quadratic penalty to the objective, then updates the multiplier from the remaining violation. It reaches exact feasibility with a moderate, fixed penalty weight, avoiding the ill-conditioning of pure penalty methods. It underlies solvers for constrained placement and sizing formulations and is the building block of ADMM.

## Minimizing the augmented Lagrangian

For minimize $f(x)$ subject to $h(x) = 0$, the augmented Lagrangian is
$$L_\rho(x, \lambda) = f(x) + \lambda h(x) + \tfrac{\rho}{2} h(x)^2.$$
Take a pin pulled toward $a$ that must sit on a track at $c$: $f = (x - a)^2$, $h = x - c$. Setting the derivative $2(x - a) + \lambda + \rho(x - c)$ to zero gives
$$x = \frac{2a - \lambda + \rho c}{2 + \rho}.$$
With $a = 4$, $c = 8$, $\rho = 1$ and $\lambda = 4$: $x = (8 - 4 + 8)/3 = 4$. The multiplier term shifts the minimizer, and once $\lambda$ equals its optimal value the minimizer sits exactly on the constraint, whatever $\rho$ is. The penalty term supplies curvature and pulls toward feasibility while $\lambda$ is still wrong.

## The multiplier update

After the minimization, update $\lambda \leftarrow \lambda + \rho\, h(x)$. This is a gradient-ascent step on the dual with step $\rho$, and it moves $\lambda$ toward the value that makes the minimizer feasible.

Example: $a = 2$, $c = 10$, $\rho = 5$, $\lambda = 6$. Then $x = (4 - 6 + 50)/7 = 6.857$, the residual is $6.857 - 10 = -3.143$, and $\lambda \leftarrow 6 + 5 \cdot (-3.143) = -9.714$. The exact multiplier is $2(a - c) = -16$: at $x = c$, stationarity of $f + \lambda h$ needs $2(c - a) + \lambda = 0$. For an equality constraint the multiplier may be negative; its sign says which way the constraint pushes. Each update reuses the previous solution as a warm start, so the inner minimizations get cheaper as the method converges.

## Convergence without a huge penalty

Two full iterations from $\lambda_0 = 0$ with $a = 13$, $c = 14$ and $\rho = 6$: $x_1 = (26 + 84)/8 = 13.75$ and $\lambda_1 = 6 \cdot (-0.25) = -1.5$; then $x_2 = (26 + 1.5 + 84)/8 = 13.9375$ and $\lambda_2 = -1.5 + 6 \cdot (-0.0625) = -1.875$. The exact multiplier is -2, and the residual shrinks by a factor of 4 each iteration, set by $2/(2 + \rho)$ here.

A pure penalty method would need $\rho \to \infty$ to reach feasibility. The augmented Lagrangian reaches it with $\rho$ fixed at a moderate value, so the inner problems stay well conditioned. Practical codes still raise $\rho$ when the violation stalls, and they apply the same idea to inequalities by penalizing only the violated part. The same iteration, split across blocks of variables, becomes ADMM.

## Going deeper

Hestenes, "Multiplier and gradient methods" (Journal of Optimization Theory and Applications, 1969); Powell, "A method for nonlinear constraints in minimization problems" (1969); Nocedal and Wright, "Numerical Optimization" (2nd ed., 2006), chapter 17.
