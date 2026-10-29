# GMRES

MNA matrices with controlled sources, inductor branches and voltage sources are unsymmetric, so CG does not apply. GMRES is the general Krylov method: it builds an orthonormal basis of the Krylov space and picks the iterate with the smallest residual. It underlies harmonic-balance solvers, preconditioned iterative SPICE variants and BEM field solvers.

## Minimal residual over the Krylov subspace

GMRES chooses $x_k \in x_0 + K_k(A, r_0)$, with $K_k = \text{span}\{r_0, Ar_0, \dots, A^{k-1}r_0\}$, minimizing $\|b - Ax_k\|_2$. The Arnoldi process builds an orthonormal basis $V_k$ of $K_k$ and a small $(k + 1) \times k$ Hessenberg matrix $H_k$ with $AV_k = V_{k+1}H_k$, so the minimization reduces to a tiny least-squares problem in $k$ unknowns, solved with Givens rotations.

Because each step minimizes over a larger space, the residual norm never increases, and in exact arithmetic GMRES terminates in at most $n$ steps. Unlike CG, there is no short recurrence for general matrices: every new basis vector must be orthogonalized against all previous ones.

## The first GMRES step by hand

From $x_0 = 0$, $K_1 = \text{span}\{b\}$, so $x_1 = \alpha b$ and the residual is $b - \alpha Ab$. Minimizing its norm is one-variable least squares:
$$\alpha = \frac{b^TAb}{\|Ab\|^2}, \qquad \|r_1\|^2 = \|b\|^2 - \frac{(b^TAb)^2}{\|Ab\|^2}.$$
For $A = \begin{pmatrix} 4 & 1 \\ 2 & 3 \end{pmatrix}$ and $b = (1, 1)$: $Ab = (5, 5)$, $b^TAb = 10$, $\|Ab\|^2 = 50$, so $\alpha = 0.2$ and $r_1 = (1, 1) - (1, 1) = 0$. Here $b$ happens to be an eigenvector, and GMRES finishes in one step.

With $b = (1, 0)$ instead: $Ab = (4, 2)$, $\alpha = 4/20 = 0.2$, residual $(0.2, -0.4)$, norm 0.447, down from 1.

## Restarts and memory

GMRES stores all $k$ basis vectors, $8kn$ bytes in double precision, and step $j$ costs about $4nj$ flops of orthogonalization on top of the matrix-vector product. With $n = 10^6$, 200 iterations need 1.6 GB for the basis alone.

Restarted GMRES($m$) runs $m$ steps, forms $x_m$, and restarts from it with a fresh basis, bounding memory at $(m + 1)n$ doubles: $m = 30$ and $n = 10^6$ is about 248 MB. The price is that restarting discards information, and GMRES($m$) can stagnate where full GMRES would converge. In practice the remedy is a good preconditioner, so that convergence happens within one restart cycle; ILU-type preconditioners and block preconditioners built from the circuit structure are common in circuit and harmonic-balance solvers.

## Going deeper

Saad and Schultz, "GMRES: a generalized minimal residual algorithm for solving nonsymmetric linear systems", SIAM J. Sci. Stat. Comput. (1986); Saad, "Iterative Methods for Sparse Linear Systems" (2nd ed., SIAM, 2003), chapter 6.
