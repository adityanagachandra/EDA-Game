# Multigrid for power grids

Multigrid solves grid problems in work proportional to the number of unknowns, with an iteration count that does not grow with grid size. It combines a cheap smoother, which removes oscillatory error, with coarse-grid corrections, which remove smooth error. For power grids with tens of millions of nodes, multigrid and algebraic multigrid preconditioners are the methods that scale.

## Smoothers kill high-frequency error

On the 1D Laplacian with $n$ interior points, the error modes are sines with frequency $\theta_k = k\pi/(n + 1)$. Damped Jacobi with weight $\omega$ multiplies mode $k$ by
$$1 - \omega(1 - \cos\theta_k)$$
per sweep. With $\omega = 2/3$, every high-frequency mode, $\theta \ge \pi/2$, is multiplied by a factor of magnitude at most $1/3$: two sweeps reduce it ninefold. A low mode, $\theta = \pi/32$, is multiplied by $1 - \frac23(1 - 0.9952) = 0.9968$, almost unchanged.

So a few relaxation sweeps leave an error that is smooth. This is the same slow convergence seen in plain Jacobi, but recast: relaxation is not a solver, it is a smoother.

## Coarse-grid correction and the V-cycle

Smooth error can be represented on a coarser grid, where it looks oscillatory again. A two-grid cycle: smooth on the fine grid; compute the residual; restrict it to the coarse grid; solve the coarse error equation; interpolate (prolong) the correction back and add it; smooth again. Applying the idea recursively gives the V-cycle.

Coarsening by 2 in each direction of a 2D grid quarters the unknowns per level. A $127 \times 127$ grid has levels of $127^2, 63^2, 31^2, \dots, 1$ unknowns, $16129 + 3969 + 961 + 225 + 49 + 9 + 1 = 21343$ in all, under $4/3$ of the finest. So a V-cycle costs a small constant times one fine-grid sweep.

## Grid-independent convergence on large power grids

A well-designed V-cycle reduces the error by a fixed factor, often 0.1 to 0.2, regardless of grid size. With factor 0.15 per cycle, reducing the error by $10^6$ takes $\lceil 6\ln 10/(-\ln 0.15)\rceil = 8$ cycles. With 3 smoothing sweeps per level and the $4/3$ hierarchy factor, that is $8 \cdot 3 \cdot 4/3 = 32$ fine-grid work units, the same on a $10^3$ or a $10^4$ grid. CG's iteration count, by contrast, grows linearly with grid width.

Real power grids are not uniform meshes: they have multiple metal layers, vias and irregular pad locations. Algebraic multigrid (AMG) builds the coarse levels from the matrix entries themselves, choosing coarse nodes from strong connections, and is often used as a preconditioner for CG. Geometric multigrid exploits the layer structure when the tool knows it.

## Going deeper

Briggs, Henson and McCormick, "A Multigrid Tutorial" (2nd ed., SIAM, 2000); Trottenberg, Oosterlee and Schüller, "Multigrid" (Academic Press, 2001).
