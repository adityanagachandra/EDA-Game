# Boundary element methods for extraction

Capacitance extraction solves Laplace's equation in the dielectric around conductors, an unbounded region. Boundary element methods put unknowns only on conductor surfaces and carry the field through space with a Green's function. The result is a small but dense system, and the fast algorithms that make it tractable are the core of field-solver extractors.

## Discretizing only the conductor surfaces

In a uniform dielectric the potential produced by a surface charge density $\sigma$ is
$$\phi(\mathbf{x}) = \int_S \frac{\sigma(\mathbf{y})}{4\pi\varepsilon|\mathbf{x} - \mathbf{y}|}\,dS_y.$$
Divide each conductor surface into $N$ panels with constant charge, and require the potential at each panel's centre to equal its conductor's voltage. That gives $Pq = v$, with $P_{ij}$ the potential at panel $i$ due to unit charge on panel $j$. Setting one conductor to 1 V and the others to 0 and summing charges gives a column of the capacitance matrix.

A box conductor meshed $10 \times 6$ on top and bottom, $10 \times 2$ on front and back and $6 \times 2$ on its ends has $2(60 + 20 + 12) = 184$ panels. No mesh of the surrounding space is needed, and the open domain is handled exactly.

## Dense matrices and their cost

Every panel interacts with every other, so $P$ is dense: $8N^2$ bytes and about $2N^3/3$ flops for LU. With $N = 20{,}000$, storage is 3.2 GB and LU about $5 \times 10^{12}$ flops. GMRES needs a matrix-vector product per iteration, $2N^2$ flops, and converges in tens of iterations because $P$ is well conditioned after simple preconditioning; with 30 iterations, LU costs about $N/90 \approx 220$ times as much.

So iterative solution wins, but even $N^2$ storage and work per product become prohibitive at hundreds of thousands of panels. Mesh refinement is costly in BEM too: doubling panel density on a surface multiplies $N$ by 4, dense storage by 16 and dense LU by 64.

## Fast multipole and hierarchical compression

Interactions between well-separated groups of panels are smooth, so they can be approximated. The fast multipole method, used in FastCap (Nabors and White, 1991), aggregates distant charges into multipole expansions and evaluates the product $Pq$ in about $O(N)$ work to a set accuracy, without forming $P$. Precorrected-FFT methods project charges onto a uniform grid and use FFTs; hierarchical matrices and adaptive cross approximation store far-field blocks in low-rank form.

Practical extractors add more: layered-media Green's functions or interface panels for stacked dielectrics, and calibration against measured or reference structures. Floating random-walk methods are an alternative family that estimate capacitance statistically with little memory. Rule-based extractors use these field solvers to build pattern libraries, then apply the library at full-chip scale.

## Going deeper

Nabors and White, "FastCap: a multipole accelerated 3-D capacitance extraction program", IEEE Trans. CAD (1991); Phillips and White, "A precorrected-FFT method for electrostatic analysis of complicated 3-D structures", IEEE Trans. CAD (1997).
