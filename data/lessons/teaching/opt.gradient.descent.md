# Gradient descent

Gradient descent moves the variables a small step against the gradient, the direction of steepest local increase. It is the workhorse of analytical placement, where millions of cell coordinates are updated each iteration from the gradient of smoothed wirelength plus density, and of continuous sizing and many legalization refinements. The mechanics are simple; the behaviour, how fast it converges and when it diverges, is governed by the curvature of the objective.

## The gradient of a quadratic cost

For $f(x, y) = \tfrac{1}{2}(ax^2 + 2hxy + by^2) - ux - vy$ the gradient is
$$\nabla f = (ax + hy - u,\ hx + by - v).$$
This is the form a cell's quadratic wirelength takes after expanding its springs: $a$ and $b$ are summed weights, $h$ is a coupling term, and $u, v$ collect fixed-pin pulls.

Example: $f = \tfrac{1}{2}(x^2 + 2xy + 4y^2) - x - 9y$ at $(1, 4)$. The x-component is $1 + 4 - 1 = 4$ and the y-component is $1 + 16 - 9 = 8$, so $\nabla f = (4, 8)$. Each component is the net force on the cell along that axis, pointing away from where the springs want it. In a placer the gradient is assembled net by net, each net adding its partial derivatives to the pins it touches, so the cost of one gradient is linear in the number of pins.

## One descent step

The update is $x_{k+1} = x_k - \eta \nabla f(x_k)$ with step size $\eta > 0$. For $f = \tfrac{1}{2}(3x^2 - 2xy + 2y^2) + 6x - 5y$ at $(5, 1)$: the gradient is $(3 \cdot 5 - 1 + 6,\ -5 + 2 - 5) = (20, -8)$. With $\eta = 0.1$ the new point is $(5 - 2,\ 1 + 0.8) = (3, 1.8)$.

All coordinates move at once, which is why gradient descent parallelizes so well on GPUs, as in DREAMPlace. The step size is the critical choice: too small and progress crawls, too large and the iterates overshoot and diverge. A placer typically normalizes the gradient or scales the step by an estimate of the curvature, since raw gradients vary over orders of magnitude between a cell on a two-pin net and a cell on a large bus.

## Iterates on a quadratic and the contraction factor

On a one-dimensional quadratic $f(x) = \tfrac{a}{2}(x - c)^2$ the gradient is $a(x - c)$, so
$$x_{k+1} - c = (1 - \eta a)(x_k - c), \qquad x_k = c + (1 - \eta a)^k (x_0 - c).$$
The error is multiplied by $r = 1 - \eta a$ each step. With $a = 4$, $c = 3$, $x_0 = 16$ and $\eta = 0.45$: $r = -0.8$, so $x_3 = 3 + (-0.8)^3 \cdot 13 = -3.656$. The negative factor makes the iterates oscillate around the minimum while shrinking.

Convergence needs $|r| < 1$, that is $0 < \eta < 2/a$. In many dimensions each eigen-direction of the Hessian has its own factor $1 - \eta\lambda_i$, so the stiffest direction limits the step and the flattest one limits the speed. That tension, the condition number, is what step-size rules, momentum and preconditioning address.

## Going deeper

Nocedal and Wright, "Numerical Optimization" (2nd ed., 2006), chapters 2 and 3; Lin et al., "DREAMPlace: deep learning toolkit-enabled GPU acceleration for modern VLSI placement" (DAC, 2019).
