# Preconditioning

Krylov methods are only as fast as the spectrum of the matrix allows, and power-grid and thermal matrices have bad spectra. A preconditioner $M \approx A$ transforms the system into one with clustered eigenvalues while staying cheap to apply. Choosing and tuning preconditioners is where most of the engineering of iterative EDA solvers happens.

## Solving an equivalent system that converges faster

Instead of $Ax = b$, solve $M^{-1}Ax = M^{-1}b$ (left preconditioning) or a symmetric form $L^{-1}AL^{-T}y = L^{-1}b$ when $M = LL^T$, which keeps CG applicable. The solution is unchanged; only the iteration sees a different matrix. The goal is $\kappa(M^{-1}A) \ll \kappa(A)$, or tightly clustered eigenvalues, with $M^{-1}r$ cheap to compute.

Since CG iterations scale like $\sqrt\kappa$, cutting $\kappa$ from $10^6$ to $10^2$ reduces 2000 iterations to about $2000\sqrt{10^{-4}} = 20$. The extremes show the trade-off: $M = I$ costs nothing and helps nothing; $M = A$ converges in one iteration but applying $M^{-1}$ is the original problem.

## Jacobi scaling of a badly scaled matrix

The Jacobi preconditioner is $M = \text{diag}(A)$. In symmetric form it gives $D^{-1/2}AD^{-1/2}$, which has a unit diagonal. For
$$A = \begin{pmatrix} 10^4 & -60 \\ -60 & 1 \end{pmatrix}:$$
$\kappa_2(A) \approx 1.6 \times 10^4$, dominated by the scale difference. Scaling gives $\begin{pmatrix} 1 & \rho \\ \rho & 1 \end{pmatrix}$ with $\rho = -60/\sqrt{10^4} = -0.6$, eigenvalues $1 \pm 0.6$, so $\kappa = 1.6/0.4 = 4$.

Jacobi removes bad scaling, from mixed units, wide conductance ranges or nonuniform meshes, but it cannot fix the $h^{-2}$ growth of a uniform grid, whose diagonal is already uniform. That needs preconditioners that capture coupling: incomplete Cholesky, or multigrid.

## Cost per iteration against iterations saved

Total time is setup plus iterations times (SpMV plus preconditioner application). Example: unpreconditioned CG needs 2000 iterations at 2 ms per SpMV, 4 s. An incomplete Cholesky preconditioner takes 150 ms to build and 3 ms per application and cuts the count to 60: $0.15 + 60 \cdot 5$ ms $= 0.45$ s, a speedup of about 8.9.

A preconditioner that halves iterations but triples the per-iteration cost loses. Setup matters when it cannot be reused; in a transient IR-drop analysis, where the matrix changes little between timesteps, one setup can serve many solves, which shifts the balance toward stronger preconditioners. Parallelism matters too: triangular solves in incomplete factorizations are sequential, while Jacobi and polynomial preconditioners map well to GPUs.

## Going deeper

Saad, "Iterative Methods for Sparse Linear Systems" (2nd ed., SIAM, 2003), chapters 9 and 10; Benzi, "Preconditioning techniques for large linear systems: a survey", J. Comput. Phys. (2002).
