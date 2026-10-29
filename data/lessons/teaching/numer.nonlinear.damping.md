# Damped Newton and step limiting

Newton's method converges only near the solution. Far from it, the full step can overshoot, oscillate or diverge, especially with exponential devices. Damped Newton keeps the Newton direction but takes a fraction of the step, chosen so the residual decreases or the voltages stay in a trusted range. It is the first robustness layer in every DC solver.

## Scaling the Newton step

Damped Newton takes $x_{k+1} = x_k + \lambda\Delta$ with $0 < \lambda \le 1$, where $\Delta$ is the full Newton step. The direction is unchanged; only the length shrinks.

Example: solve $x^3 - 8 = 0$ from $x_0 = 0.5$. The full step is $\Delta = -(0.125 - 8)/(3 \cdot 0.25) = 10.5$, which would land at $x = 11$, far past the root 2. With $\lambda = 0.25$, $x_1 = 0.5 + 2.625 = 3.125$, much closer.

The Newton direction is always a descent direction for the merit function $\frac12\|F\|^2$ when $J$ is nonsingular, so a small enough $\lambda$ always reduces the residual. Damping trades speed for safety: near the root $\lambda = 1$ must be restored, or the quadratic convergence is lost.

## Backtracking until the residual drops

A line search chooses $\lambda$ automatically. Simple backtracking tries $\lambda = 1, 1/2, 1/4, \dots$ and accepts the first with $\|F(x_k + \lambda\Delta)\| < \|F(x_k)\|$; the Armijo rule asks for a sufficient decrease, $\|F\| \le (1 - c\lambda)\|F(x_k)\|$ with a small $c$ such as $10^{-4}$.

The classic test case is $f(x) = \arctan x$: from $|x_0| > 1.39$ plain Newton diverges, each step overshooting further. From $x_0 = 2$: $\Delta = -\arctan(2)(1 + 4) = -5.536$. The full step lands at -3.536, where $|\arctan| = 1.295 > 1.107$, rejected. With $\lambda = 1/2$ the trial point is -0.768, $|\arctan| = 0.655 < 1.107$, accepted. The next iterations converge to 0.

## Limiting the voltage change per iteration

Circuit simulators often bound the change of node voltages instead: scale the whole step by $\lambda = \min(1, \Delta V_{max}/\max_i|\Delta_i|)$. With proposed updates $(0.3, -1.6, 2.0)$ V and $\Delta V_{max} = 0.5$ V, $\lambda = 0.5/2.0 = 0.25$, and the updates become $(0.075, -0.4, 0.5)$ V.

This global scaling keeps the Newton direction. SPICE's device-level limiting is different: pnjlim and fetlim adjust each device's junction or gate voltage before evaluating it, using logarithmic steps for junctions. That changes the direction and is not a line search, but it is very effective against exponential overflow. Robust solvers combine both, and fall back to continuation methods (gmin or source stepping) when damping alone fails.

## Going deeper

Dennis and Schnabel, "Numerical Methods for Unconstrained Optimization and Nonlinear Equations" (SIAM, 1996), chapter 6; Kelley, "Solving Nonlinear Equations with Newton's Method" (SIAM, 2003).
