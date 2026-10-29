# The conductance matrix of a power grid

IR-drop analysis writes the grid as a linear system $G d = i$, where $d$ holds each node's drop below the supply and $i$ its load current. The matrix $G$ is built by stamping each resistor, exactly as in nodal analysis. Its structure, sparse, symmetric and positive definite, decides which solvers work and how large a grid a tool can handle. A PDN engine designer must be able to write $G$ by hand for a small grid.

## Diagonal entries

Each branch of conductance $g = 1/R$ between nodes $j$ and $k$ adds $g$ to $G_{jj}$ and $G_{kk}$. A supply connection of conductance $g_p$ at node $j$ adds $g_p$ to $G_{jj}$ only, because the supply side is fixed at zero drop. So a diagonal entry is the sum of all conductances touching the node. In a 2 by 2 grid where n4 connects to n2 through $2\ \Omega$, to n3 through $1\ \Omega$ and to the supply through $1\ \Omega$, $G_{44} = 0.5 + 1 + 1 = 2.5$ S.

## Off-diagonal entries and sparsity

The same branch adds $-g$ to $G_{jk}$ and $G_{kj}$, so $G$ is symmetric. A branch of $4\ \Omega$ between n1 and n3 gives $G_{13} = -0.25$ S. Nodes with no branch between them, such as diagonal neighbours n1 and n4 in a grid, have a zero entry. In a real grid each node touches only a few neighbours, so a row has about five non-zeros out of millions of columns. The tool stores $G$ in a sparse format and never forms dense rows; fill-in during factorisation is what limits direct solvers.

## Row sums and definiteness

Every internal branch contributes $+g$ on a diagonal and $-g$ off it in the same row, so it cancels in the row sum. Only supply connections survive: the row sum at a node equals its supply conductance, and it is zero for a node with no supply connection. With no supply anywhere, every row sums to zero, the all-ones vector is in the null space and $G$ is singular, the graph Laplacian of the grid. One supply connection anywhere in a connected grid makes $G$ symmetric positive definite. That is the property conjugate gradient and Cholesky factorisation need.

## Going deeper

Najm, "Circuit Simulation" (2010), the chapter on nodal analysis; Kahng et al., "VLSI Physical Design" (2nd ed., 2022).
