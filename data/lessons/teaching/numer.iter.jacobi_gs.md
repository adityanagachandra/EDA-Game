# Jacobi and Gauss-Seidel on power grids

The simplest iterative solvers update one node at a time from its neighbours. On a power grid each update is a local KCL balance, so Jacobi and Gauss-Seidel are easy to implement, trivially parallel (Jacobi) and memory-light. They are also very slow on large grids, and understanding why sets up conjugate gradients, preconditioning and multigrid.

## A node update averages its neighbours

At a grid node $i$ with conductances $g_j$ to neighbours at voltages $v_j$ and a load drawing current $I_i$, KCL reads $\sum_j g_j(v_j - v_i) = I_i$. Solving for $v_i$ with the neighbours held fixed gives
$$v_i = \frac{\sum_j g_j v_j - I_i}{\sum_j g_j}.$$
A node tied by 10 S to a pad at 1.000 V and by 5 S to a neighbour at 0.970 V, with a 0.3 A load: $v_i = (10 + 4.85 - 0.3)/15 = 0.9700$ V.

The update is a conductance-weighted average of neighbours, shifted down by the load. Jacobi applies it to every node using only old values. Gauss-Seidel sweeps through the nodes and uses each new value immediately, which typically converges about twice as fast but makes the sweep sequential unless nodes are coloured (red-black ordering on a grid).

## One sweep on a resistive grid

Take a rail: pad at 1.0 V, nodes 1, 2, 3, pad at 1.0 V, 10 S between neighbours, loads 0.1 A at each node, all nodes starting at 1.0 V. With equal conductances each update is $v_i = (v_{left} + v_{right} - I_i/g)/2$.

Jacobi (old values): every node gets $(1 + 1 - 0.01)/2 = 0.995$ V.
Gauss-Seidel in order 1, 2, 3: $v_1 = 0.995$; $v_2 = (0.995 + 1 - 0.01)/2 = 0.9925$; $v_3 = (0.9925 + 1 - 0.01)/2 = 0.99125$ V.

The exact answer has $v_2 = 0.98$ V at the middle; both methods are still far from it after one sweep, because information travels only one node per sweep in Jacobi. On a grid of width $k$, a load in the middle needs on the order of $k$ sweeps just to feel the pads, and many more to converge.

## Convergence rate and diagonal dominance

For a chain of $n$ interior nodes between two fixed pads with equal conductances, the Jacobi iteration matrix has spectral radius $\rho = \cos(\pi/(n + 1))$, and Gauss-Seidel has $\rho^2$. Reducing the error by $10^{-d}$ takes about $d\ln 10/(-\ln\rho)$ iterations. For $n = 100$: $\rho = 0.99952$, and $d = 6$ needs about 28,600 Jacobi iterations; Gauss-Seidel needs half as many.

Since $1 - \cos(\pi/(n+1)) \approx \pi^2/(2(n+1)^2)$, the iteration count grows like $n^2$. Convergence is guaranteed because grounded resistive grids give irreducibly diagonally dominant M-matrices, but slowly: the smooth, long-wavelength error is what these methods fail to remove. That is exactly what multigrid's coarse grids target, and why IR-drop tools use preconditioned CG instead of plain relaxation.

## Going deeper

Saad, "Iterative Methods for Sparse Linear Systems" (2nd ed., SIAM, 2003), chapter 4; Briggs, Henson and McCormick, "A Multigrid Tutorial" (2nd ed., SIAM, 2000), chapter 2.
