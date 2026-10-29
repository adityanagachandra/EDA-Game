# Solving the placement system by conjugate gradient

Every iteration of a quadratic placer solves $Ax = b$ for millions of cells. Direct factorization is too slow and memory-hungry; conjugate gradient (CG), which needs only matrix-vector products with the sparse A, is the standard solver. Its first step, its convergence rate and its sensitivity to conditioning explain both how fast a quadratic placer runs and why preconditioners and warm starts matter.

## The first conjugate gradient step

CG starts from a guess $x_0$, here 0, with residual $r_0 = b - Ax_0 = b$ and search direction $p_0 = r_0$. The step length minimises the quadratic along $p_0$:
$$\alpha_0 = \frac{r_0^T r_0}{p_0^T A p_0}, \qquad x_1 = x_0 + \alpha_0 p_0.$$
For $A = \begin{pmatrix} 9 & -3 \\ -3 & 9 \end{pmatrix}$ and $b = (39, 12)$: $r_0^T r_0 = 1665$, $Ap_0 = (315, -9)$, $p_0^T A p_0 = 12177$, so $\alpha_0 = 0.1367$ and $x_1 = (5.33, 1.64)$. Later steps update $r_{k+1} = r_k - \alpha_k A p_k$ and choose the next direction A-conjugate to the previous ones. The only expensive operation in each step is one sparse matrix-vector product, $Ap_k$, whose cost is proportional to the number of nonzeros, that is, to the number of two-pin connections in the net model. Everything else is a handful of dot products and vector updates, which is why the method scales to millions of cells and parallelises well.

## Convergence and conditioning

In exact arithmetic CG reaches the solution in at most $n$ steps, but for large systems what matters is the rate: the error in the A-norm obeys $\|e_k\|_A \le 2\rho^k\|e_0\|_A$ with $\rho = (\sqrt{\kappa} - 1)/(\sqrt{\kappa} + 1)$, where $\kappa = \lambda_{max}/\lambda_{min}$ is the condition number. For the matrix above the eigenvalues are 12 and 6, $\kappa = 2$ and $\rho = 0.1716$, fast. Placement matrices are worse: cells far from any fixed pin make $\lambda_{min}$ small, and high-degree cells make $\lambda_{max}$ large. The number of iterations for a given accuracy grows roughly like $\sqrt{\kappa}$.

## Making the solve fast in a placer

Placers use three tricks. Preconditioning: solving $M^{-1}Ax = M^{-1}b$ with a cheap approximation M of A, such as its diagonal (Jacobi) or an incomplete Cholesky factor, shrinks the effective condition number. Warm starts: each placement iteration changes A and b only a little, so starting CG from the previous solution leaves a small initial error. Loose tolerances: the solution is followed by spreading, so solving to high precision is wasted effort, and a few dozen iterations usually suffice. The x and y systems are independent and can be solved in parallel, and the sparse matrix-vector product is the only heavy kernel, which maps well to multicore CPUs and GPUs.

## Going deeper

Shewchuk, "An introduction to the conjugate gradient method without the agonizing pain" (1994); Saad, "Iterative Methods for Sparse Linear Systems" (2nd ed., 2003), chapters 6 and 10.
