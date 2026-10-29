# Sparsity and structure of the MNA matrix

Circuit matrices are extremely sparse: a node connects to a handful of neighbours, so a row of a million-node matrix holds perhaps five nonzeros. The pattern of nonzeros is fixed by the netlist before any value is known, and simulators analyze it once, choose an elimination order, and reuse the symbolic factorization for every Newton iteration and time step. Counting nonzeros and spotting zero diagonals is how a tool designer reasons about memory, fill-in and pivoting.

## Counting structural nonzeros

A structural nonzero is a position that some stamp touches, whatever its value. In the nodal conductance matrix every node has a nonzero diagonal, and each distinct pair of non-ground nodes joined by at least one resistor adds two symmetric off-diagonal positions. Resistors to ground touch only diagonals, and parallel resistors share positions. So the count is $n + 2p$, where $p$ is the number of distinct connected node pairs. A three-node chain 1-2-3 with a resistor to ground has $3 + 2 \cdot 2 = 7$ nonzeros out of 9. Simulators store only these entries, in compressed column or row form, with the pattern computed once from the netlist.

## Nonzeros of the augmented matrix

Voltage sources add a branch row and column. A source from node $a$ to ground adds $+1$ at $(a, k)$ and $(k, a)$, where $k$ is its branch index: 2 nonzeros. A source between two non-ground nodes adds 4. Its own diagonal $(k, k)$ stays empty. A five-node network with 13 nonzeros in its node block plus one grounded voltage source has $13 + 2 = 15$ nonzeros in a 6 by 6 matrix. Because stamps never change position between iterations, the simulator can map every element to the matrix addresses it writes during setup, so the load phase is a sequence of direct indexed additions.

## Zero diagonals and pivoting

Every branch row that expresses a voltage constraint (independent voltage source, VCVS, or inductor at DC) has no term in its own current, so its diagonal is a structural zero. The count of zero diagonals is the number of such branches. Node rows have positive diagonals from conductances and, in practice, from the gmin added across junctions. Zero diagonals mean the matrix cannot be factored in an arbitrary order without pivoting. SPICE-class solvers choose a pivot order that balances numerical stability against fill-in, the new nonzeros created by elimination, often by ordering for sparsity and then checking pivot magnitudes. A bad order can turn a sparse matrix dense.

## Going deeper

Davis, "Direct Methods for Sparse Linear Systems" (SIAM, 2006); Davis and Palamadai Natarajan, "Algorithm 907: KLU, a direct sparse solver for circuit simulation problems" (ACM TOMS, 2010).
