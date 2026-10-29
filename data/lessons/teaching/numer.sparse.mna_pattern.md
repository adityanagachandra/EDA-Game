# Nonzero pattern of a nodal matrix

Every circuit simulator, IR-drop analyzer and extractor ends up solving a linear system whose matrix comes from a netlist. Before any number is computed, the netlist already fixes which entries of that matrix can be nonzero. That pattern decides how the matrix is stored, how much fill elimination creates and how fast an iterative solver runs, so a tool designer reads the pattern straight off the circuit graph.

## One unknown per node, one entry per connection

Nodal analysis writes Kirchhoff's current law at every non-ground node. A resistor of conductance $g$ between nodes $i$ and $j$ adds $g$ to $G_{ii}$ and $G_{jj}$ and $-g$ to $G_{ij}$ and $G_{ji}$. A resistor from node $i$ to ground only adds $g$ to $G_{ii}$, because ground is not an unknown. Two resistors in parallel between the same nodes land on the same four entries and simply add.

So the off-diagonal pattern of $G$ is the adjacency matrix of the circuit graph with ground deleted, and the diagonal is full. Row $i$ has $1 + \deg(i)$ nonzeros, where $\deg(i)$ counts distinct non-ground neighbours. Take node 1 joined to nodes 2 and 4 and to ground: row 1 has 3 nonzeros, the diagonal and columns 2 and 4. The matrix is the weighted graph Laplacian plus a diagonal of grounded conductances.

## Counting nonzeros from the netlist

For a network of resistors on $n$ non-ground nodes, let $m$ be the number of distinct unordered pairs of non-ground nodes joined by at least one element. Then
$$\text{nnz}(G) = n + 2m.$$
Elements to ground contribute nothing new, and parallel elements are counted once. Example: 5 nodes; resistors 1-2, 2-3, 3-4, 4-5, 1-5 and a second resistor 2-3 in parallel; resistors from nodes 1 and 4 to ground. The distinct pairs are 1-2, 2-3, 3-4, 4-5, 1-5, so $m = 5$ and $\text{nnz} = 5 + 10 = 15$.

Modified nodal analysis adds a row and column for each voltage source or inductor branch current; those add a few unsymmetric entries, but the counting idea is the same: one entry per (equation, unknown) coupling that some element stamps.

## Why circuit matrices are so sparse

A node in a real circuit touches a handful of devices, so the average row holds a small constant number of entries while $n$ grows into the millions. A $k$-by-$k$ power-grid mesh has $n = k^2$ nodes and $2k(k - 1)$ resistors, so
$$\text{nnz} = k^2 + 4k(k - 1) \approx 5n.$$
For $k = 1000$ that is about $5 \times 10^6$ stored entries against $10^{12}$ for a dense matrix, a density near $5 \times 10^{-6}$.

Sparsity is what makes circuit-scale linear algebra possible at all: storage must be linear in nnz, matrix-vector products cost $2\,\text{nnz}$ flops, and direct factorization is fast only if elimination keeps the factors sparse too. Supply nets are the exception: a VDD node connected to thousands of devices gives a dense row and column, which is why orderings treat such hubs with care.

## Going deeper

Vlach and Singhal, "Computer Methods for Circuit Analysis and Design" (2nd ed.), chapters on nodal and modified nodal formulation; Davis, "Direct Methods for Sparse Linear Systems" (SIAM, 2006), chapter 1.
