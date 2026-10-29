# Quasi-Newton methods: BFGS and L-BFGS

Quasi-Newton methods build a Hessian approximation from successive gradients instead of computing second derivatives. BFGS updates a dense approximation with a rank-two correction each step and converges superlinearly; L-BFGS keeps only the last few step and gradient-change pairs, so its memory is linear in the number of variables. L-BFGS is a standard choice for smooth large-scale problems, including analytical placement and sizing formulations, when a plain gradient is too slow and a true Hessian too costly.

## The secant condition in one dimension

A quasi-Newton approximation B must satisfy the secant condition $B s = y$, where $s = x_{k+1} - x_k$ is the step and $y = \nabla f(x_{k+1}) - \nabla f(x_k)$ the gradient change. In one dimension that fixes B completely: $B = y/s$, the slope of the derivative between the two points, and the update $x_{k+2} = x_{k+1} - f'(x_{k+1})/B$ is the secant method applied to $f'$.

Example: $f'(2) = -11$ and $f'(12) = -5$. Then $B = 6/10 = 0.6$ and $x_2 = 12 - (-5)/0.6 = 20.333$. The model predicts that the derivative reaches zero 8.33 units beyond 12. If the curvature estimate were negative, the step would point uphill, which is why quasi-Newton methods require $y^T s > 0$ and enforce it with a line search satisfying the Wolfe conditions.

## The BFGS update

In several dimensions the secant condition does not determine B, so BFGS picks the update closest to the previous matrix that satisfies it, a rank-two correction:
$$B_{k+1} = B_k - \frac{B_k s s^T B_k}{s^T B_k s} + \frac{y y^T}{y^T s}.$$
From $B_0 = I$ with $s = (-3, -3)$ and $y = (1, -4)$: $s^T s = 18$, $y^T s = -3 + 12 = 9 > 0$. Entry $(2, 2)$ is $1 - 9/18 + 16/9 = 2.278$, and entry $(1, 2)$ is $0 - 9/18 + (1)(-4)/9 = -0.944$.

If $B_k$ is positive definite and $y^T s > 0$, then $B_{k+1}$ is positive definite, so every quasi-Newton step is a descent direction. In practice implementations update the inverse approximation directly, which avoids solving a linear system each step.

## Limited memory at placement scale

Dense BFGS stores an n by n matrix: impossible for $n = 4 \times 10^6$ coordinates of two million cells, which would need about $1.3 \times 10^{14}$ bytes. L-BFGS stores only the last m pairs $(s_k, y_k)$, typically $m$ between 3 and 20, and applies the inverse-Hessian approximation to a vector with a two-loop recursion in $O(mn)$ time.

Memory is $2mn$ numbers. With $n = 4 \times 10^6$, $m = 13$ and 8-byte doubles, that is $2 \cdot 13 \cdot 4 \times 10^6 \cdot 8 = 832$ MB, large but feasible, and $m = 5$ needs 320 MB. L-BFGS works best on smooth objectives. Placement objectives change shape as density weights rise, which makes stored curvature pairs stale, so placers often reset the memory when the penalty weight changes or prefer Nesterov's method, which needs no curvature history.

## Going deeper

Nocedal, "Updating quasi-Newton matrices with limited storage" (Mathematics of Computation, 1980); Liu and Nocedal, "On the limited memory BFGS method for large scale optimization" (Mathematical Programming, 1989); Nocedal and Wright, "Numerical Optimization" (2nd ed., 2006), chapters 6 and 7.
