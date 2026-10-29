# The Maxwell capacitance matrix

A field solver's output for $n$ conductors is the Maxwell capacitance matrix, which relates their charges to their voltages. Circuit netlists instead use ground and coupling capacitors. Converting between the two forms, and knowing the matrix's physical properties, is the interface between extraction and simulation, and those properties double as cheap correctness checks on any solver.

## Building the Maxwell matrix

The Maxwell matrix $C$ satisfies $Q = CV$. Its diagonal entry $C_{ii}$ is the total capacitance of conductor $i$: its capacitance to ground plus all its coupling capacitances. Its off-diagonal entry $C_{ij}$ is minus the coupling capacitance between $i$ and $j$. So from a SPICE-style description with ground capacitances $C_{g,i}$ and couplings $C_{c,ij}$:
$$C_{ii} = C_{g,i} + \sum_{j \ne i} C_{c,ij}, \qquad C_{ij} = -C_{c,ij}.$$

Example: three wires with ground capacitances 10, 8 and 12 fF and couplings $C_{12} = 4$, $C_{13} = 1$, $C_{23} = 5$ fF. Then $C_{11} = 15$, $C_{22} = 17$, $C_{33} = 18$, $C_{12} = -4$, $C_{13} = -1$ and $C_{23} = -5$ fF. The matrix is symmetric, which is reciprocity, and diagonally dominant, because each diagonal exceeds the sum of the magnitudes of its row's couplings by exactly the ground capacitance.

## Ground capacitance from row sums

Raising every conductor to 1 V leaves no voltage across any coupling capacitor, so only the ground capacitances hold charge. Hence the ground capacitance of conductor $i$ is its row sum:
$$C_{g,i} = \sum_j C_{ij}.$$
Couplings are read directly as $C_{c,ij} = -C_{ij}$.

Example: a solver reports the row $(15, -4, -1)$ fF for conductor 1, so its ground capacitance is 10 fF and its couplings are 4 fF and 1 fF. For a passive structure every row sum is non-negative and every off-diagonal entry is non-positive. A negative row sum is unphysical and points to an unconverged solve or a mesh problem, and a significant asymmetry $C_{ij} \ne C_{ji}$ points to the same. These checks cost nothing and catch real solver bugs.

## Charge from the matrix

With all voltages known, each conductor's charge is $Q_i = \sum_j C_{ij}V_j$. Femtofarads times volts give femtocoulombs.

Example, continuing with $C_{11} = 15$, $C_{12} = -4$ and $C_{13} = -1$ fF: at voltages 1.0, 0.5 and 0 V, $Q_1 = 15 - 2 - 0 = 13$ fC. Equivalently, the ground capacitor holds $10 \cdot 1.0 = 10$ fC and the couplings hold $4(1.0 - 0.5) + 1(1.0 - 0) = 3$ fC, the same 13 fC. The matrix form makes the dependence on neighbours' voltages explicit, which is the root of the Miller effect: the charge a driver must deliver depends on what its neighbours do. The matrix is also symmetric positive definite, so the stored energy $\frac12 V^TCV$ is positive for any nonzero voltages.

## Going deeper

Paul, "Analysis of Multiconductor Transmission Lines" (2nd ed.), chapter 3; Nabors and White, "FastCap: a multipole accelerated 3-D capacitance extraction program", 1991.
