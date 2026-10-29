# The method of moments field solver

The method of moments (MoM), or boundary-element method, solves electrostatics by placing unknown charges on conductor surfaces only, using the free-space Green's function for the space between. It is the classic engine of 3D capacitance extraction, implemented in tools such as FastCap. Its strength is a small surface mesh; its cost is a dense matrix, which fast algorithms then compress.

## Panels and the potential coefficient

Discretize every conductor surface into panels, each carrying an unknown uniform charge $q_j$. The potential at panel $i$ is the sum of contributions from all panels, $\phi_i = \sum_j P_{ij}q_j$, where $P_{ij}$ integrates the Green's function $1/(4\pi\varepsilon_0 r)$ over panel $j$ as seen from panel $i$'s centre. For panels far apart compared with their size, $P_{ij} \approx 1/(4\pi\varepsilon_0 r_{ij})$.

Numerically, $1/(4\pi\varepsilon_0) = 8.988 \times 10^9$ m/F, which is 8.988 V per pC at 1 mm. Example: charges of 0.5 pC at 1 mm and -0.2 pC at 2 mm produce $8.988(0.5/1 - 0.2/2) = 3.60$ V at a third panel. The self term $P_{ii}$ requires integrating a singular kernel over the panel itself, which is the delicate part of any MoM implementation.

## Solving for the charges

Setting each panel's potential to its conductor's voltage gives the dense system $Pq = v$. Solving it with conductor $k$ at 1 V and all others at 0 V yields the charges that form column $k$ of the Maxwell capacitance matrix: the total charge on conductor $i$ is $C_{ik}$.

Example with one panel per conductor and $P = \begin{pmatrix} a & b \\ b & a \end{pmatrix}$ in V/pC: $P^{-1} = \frac{1}{a^2 - b^2}\begin{pmatrix} a & -b \\ -b & a \end{pmatrix}$. With $a = 20$ and $b = 5$, raising conductor 1 to 1 V gives $q_1 = 20/375 = 0.0533$ pC and $q_2 = -5/375 = -0.0133$ pC, so $C_{11} = 0.0533$ pF and $C_{21} = -0.0133$ pF. Conductor 2, grounded, acquires negative induced charge. Real models use thousands of panels per conductor, refined near edges where charge concentrates.

## Dense matrices and their cost

Every panel interacts with every other, so $P$ is dense: $N$ panels need $N^2$ entries and a direct LU solve costs about $\frac{2}{3}N^3$ operations. Example: 20,000 panels in double precision need $8 \cdot 20000^2 = 3.2 \times 10^9$ bytes, 3.2 GB, and doubling the panel count multiplies the factorization time by 8.

Production solvers avoid both costs. They solve iteratively with a Krylov method such as GMRES, which needs only matrix-vector products, and they compute those products approximately in $O(N)$ or $O(N\log N)$ time with the fast multipole method (as in FastCap), precorrected FFT, or hierarchical low-rank compression. Accuracy then depends on the compression tolerance as well as the mesh, so a solver should report both. This scaling argument is why MoM is the accurate mode and pattern tables are the fast one in full-chip extraction.

## Going deeper

Nabors and White, "FastCap: a multipole accelerated 3-D capacitance extraction program", IEEE Transactions on Computer-Aided Design, 1991; Harrington, "Field Computation by Moment Methods"; Gibson, "The Method of Moments in Electromagnetics" (2nd ed.).
