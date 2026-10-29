# The Maxwell capacitance matrix

A set of conductors over ground is described exactly, in electrostatics, by its capacitance matrix: the linear map from conductor voltages to conductor charges. Field solvers output it, extractors read ground and coupling capacitances out of it, and confusing its sign conventions is a classic integration bug. A tool designer needs to read the matrix fluently in both directions: from matrix to circuit elements, and from elements back to the matrix.

## Couplings off the diagonal

For $N$ conductors with voltages $V$ relative to ground, the charges are $Q = C V$, where $C$ is the Maxwell capacitance matrix. Entry $C_{ij}$ for $i \ne j$ is the charge induced on conductor $i$ when conductor $j$ is at 1 V and all others at 0 V. That charge is negative: raising one conductor pulls opposite charge onto its grounded neighbours. The circuit coupling capacitance between $i$ and $j$ is therefore $-C_{ij}$, a positive number. If row 2 of a three-conductor matrix is $[-8, 17, -2]$ fF, conductor 2 couples 8 fF to conductor 1 and 2 fF to conductor 3. The matrix is symmetric by reciprocity, so $C_{23} = C_{32}$ is a free sanity check on a solver run.

## Ground capacitance as a row sum

Put every conductor at 1 V. Nothing then differs in voltage between conductors, so the couplings carry no charge, and conductor $i$ holds only the charge on its capacitance to ground, $C_{gi}$. But by the matrix, its charge is $\sum_j C_{ij}$, the row sum. So the ground capacitance is the row sum. For row 3, $[-4, -2, 10]$ fF, ground capacitance is $-4 - 2 + 10 = 4$ fF. If a design has no explicit ground conductor, rows sum to zero, which signals that the reference was chosen as one of the conductors. Extraction tools convert the matrix into a netlist of grounded and coupling capacitors by exactly these two rules.

## Building the diagonal

Going the other way, the diagonal entry is the total capacitance of the conductor: its ground capacitance plus all its couplings, $C_{ii} = C_{gi} + \sum_{j \ne i} C_{c,ij}$. With ground capacitances of 2, 8 and 2 fF and couplings of 5 (1 to 2), 6 (1 to 3) and 8 fF (2 to 3), $C_{33} = 2 + 6 + 8 = 16$ fF. This is the load a driver of conductor 3 sees when every neighbour is held still. The matrix is diagonally dominant and positive definite, properties that solvers and reducers rely on. A common bug is mixing this Maxwell form with the mutual form that some tools report, where off-diagonals are positive coupling values and diagonals are ground capacitances only.

## Going deeper

Jackson, "Classical Electrodynamics" (3rd ed., 1999), the section on capacitance coefficients; Nabors and White, "FastCap: a multipole accelerated 3-D capacitance extraction program", IEEE Transactions on CAD (1991).
