# Conjugate gradient for placement systems

Conjugate gradient minimizes a convex quadratic, or equivalently solves a symmetric positive definite system $Ax = b$, using only matrix-vector products. It is how quadratic placers solve their million-variable sparse systems. Each step does an exact line search along a direction that is A-conjugate to all previous ones, so in exact arithmetic it finishes in at most n steps, and in practice a good preconditioner makes it converge in far fewer.

## The first step length

Minimizing $\tfrac{1}{2}x^T A x - b^T x$ from $x_0 = 0$, the residual is $r_0 = b - Ax_0 = b$, which is minus the gradient, and the first direction is $d_0 = r_0$. The exact line minimum along $d_0$ is at
$$\alpha_0 = \frac{r_0^T r_0}{d_0^T A d_0}.$$
Example: $A = \begin{pmatrix} 5 & 0 \\ 0 & 7 \end{pmatrix}$ and $b = (3, 4)$. Then $r_0^T r_0 = 25$, $Ad_0 = (15, 28)$, $d_0^T A d_0 = 45 + 112 = 157$, and $\alpha_0 = 25/157 = 0.1592$.

This is steepest descent with the ideal step for a quadratic, computed without any trial evaluations because the curvature along $d_0$ is known exactly. In a placer, $Ad_0$ is one sparse matrix-vector product, a pass over the nets.

## The first iterate and the next direction

The first iterate is $x_1 = x_0 + \alpha_0 d_0$. With $A = \begin{pmatrix} 8 & 1 \\ 1 & 5 \end{pmatrix}$ and $b = (-1, 10)$: $r_0^T r_0 = 101$, $Ad_0 = (2, 49)$, $d_0^T A d_0 = -2 + 490 = 488$, $\alpha_0 = 0.207$, and $x_1 = (-0.207, 2.070)$.

The method then updates the residual $r_1 = r_0 - \alpha_0 A d_0$, computes $\beta_0 = r_1^T r_1 / r_0^T r_0$, and sets the new direction $d_1 = r_1 + \beta_0 d_0$. The direction mixes the new residual with the old direction so that $d_1^T A d_0 = 0$: the directions are conjugate. Minimizing along $d_1$ then never undoes the progress made along $d_0$, which is exactly what plain steepest descent fails to guarantee and why it zigzags in narrow valleys.

## Finite termination and preconditioning

Because the directions are A-conjugate, after k steps CG has minimized the quadratic over the span of the first k directions. In exact arithmetic it therefore terminates in at most n steps for an n by n system, and in at most k steps if A has only k distinct eigenvalues. For a 2 by 2 system the second iterate is the exact solution $A^{-1}b$: with $A = \begin{pmatrix} 2 & -1 \\ -1 & 3 \end{pmatrix}$ and $b = (0, 7)$, $x_2 = (1.4, 2.8)$.

For large systems the error falls roughly by a factor $(\sqrt{\kappa} - 1)/(\sqrt{\kappa} + 1)$ per step, much better than steepest descent's $(\kappa - 1)/(\kappa + 1)$. Placers still precondition: a Jacobi (diagonal) preconditioner divides by each cell's total net weight, and incomplete Cholesky does better. Quadratic placers stop after a modest residual reduction, since the placement will be spread and re-solved anyway.

## Going deeper

Shewchuk, "An introduction to the conjugate gradient method without the agonizing pain" (technical report, Carnegie Mellon University, 1994); Saad, "Iterative Methods for Sparse Linear Systems" (2nd ed., 2003), chapter 6.
