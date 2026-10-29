# Thermal simulation

Thermal simulation solves the heat equation over the package and die to find the full temperature field, not just a few node values. The standard approaches discretize space by finite differences, finite volumes or finite elements, turning the partial differential equation into a large sparse linear system. Understanding the discretization explains the cost, the accuracy and the structure that fast solvers exploit.

## One Jacobi update

Discretize a uniform 1-D bar into nodes, each joined to its neighbours by a thermal resistance $R$. Energy balance at node $i$ with injected power $Q_i$ is $(T_{i-1} - T_i)/R + (T_{i+1} - T_i)/R + Q_i = 0$, so $T_i = (T_{i-1} + T_{i+1} + Q_i R)/2$. A Jacobi sweep updates every node from its neighbours' old values. With neighbours at 60 and 66 °C, $R = 2$ K/W and $Q_i = 1.5$ W, the update is $(60 + 66 + 3)/2 = 64.5$ °C. This is the finite-difference form of $k T'' + q''' = 0$, and it shows the method's character: each update is cheap and local, but information travels one node per sweep, so plain Jacobi converges slowly on fine grids.

## A chain with distributed heating

Exact solutions of small networks check a solver. Take $N$ nodes in a chain, the first joined to a sink at $T_s$ through resistance $R$, each joined to the next by $R$, with power $q$ injected at every node and an adiabatic far end. The link into node $k$ carries the heat of nodes $k$ to $N$, $(N - k + 1)q$, so the far-end rise is $R q \sum_{k=1}^{N} (N-k+1) = R q N(N+1)/2$. For $N = 6$, $R = 0.5$ K/W, $q = 2$ W and $T_s = 40$ °C, the far end is $40 + 0.5 \times 2 \times 21 = 61$ °C. Such closed forms are golden tests for a thermal solver.

## Grid size and sparsity

A 3-D finite-difference grid of $n_x n_y n_z$ cells gives one unknown per cell and a 7-point stencil: each row of the matrix has the node and up to six neighbours. The number of nonzeros is $N + 2[(n_x-1)n_y n_z + n_x(n_y-1)n_z + n_x n_y(n_z-1)]$, counting each neighbour link twice. A $100 \times 100 \times 10$ grid has $10^5$ unknowns and about $6.76 \times 10^5$ nonzeros. The matrix is symmetric positive definite, so conjugate gradient with a good preconditioner, or multigrid, solves it efficiently. Package-scale models reach millions of cells because thin layers need fine vertical resolution, which is why thermal tools use nonuniform meshes and model order reduction.

## Going deeper

Huang and others, "HotSpot: a compact thermal modeling methodology for early-stage VLSI design", IEEE Transactions on VLSI Systems (2006); Saad, "Iterative Methods for Sparse Linear Systems" (2nd ed., 2003).
