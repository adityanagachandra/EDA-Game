# The finite element method

The finite element method discretizes PDEs on unstructured meshes, which fit the irregular geometry of packages, TSVs and interposers far better than grids. It starts from an integral (weak) form, builds a small matrix for each element and assembles them. FEM is the method behind thermal and thermomechanical package analysis and many field solvers.

## The weak form and hat functions

For steady conduction $-(ku')' = f$ on an interval, multiply by a test function $v$ that vanishes where $u$ is prescribed and integrate by parts:
$$\int k\,u'v'\,dx = \int f\,v\,dx.$$
This weak form needs only first derivatives, so piecewise-linear functions suffice. Choose nodes $x_0 < x_1 < \dots$ and hat functions $\phi_i$, equal to 1 at node $i$, 0 at the others and linear in between. Write $u \approx \sum_j u_j\phi_j$ and take $v = \phi_i$ for each $i$: the result is a linear system $Ku = F$ with $K_{ij} = \int k\phi_i'\phi_j'\,dx$.

$K$ is symmetric positive definite once a boundary value is fixed, because $u^TKu = \int k(u')^2dx$ is the conduction energy, mirroring the resistor-network power form.

## Element stiffness matrices and assembly

On one element of length $h$, cross-section $A$ and conductivity $k$, the hat functions have slopes $\pm 1/h$, and
$$K_e = \frac{kA}{h}\begin{pmatrix} 1 & -1 \\ -1 & 1 \end{pmatrix}.$$
With $k = 150$ W/(m K), $A = 1$ mm$^2$ and $h = 2$ mm, $kA/h = 75$ mW/K. This is exactly a thermal conductance between the element's two nodes.

Assembly adds each element matrix into the global matrix at its node indices. For three elements in series with conductances 10, 20 and 30 mW/K on nodes 1 to 4, the global diagonal is $(10, 30, 50, 30)$ and the off-diagonals are $-10, -20, -30$. In 2D and 3D, triangle and tetrahedron element matrices come from the same integrals, and assembly is a scatter-add loop over elements, the FEM analogue of stamping.

## Load vectors and solving a small mesh

The right side collects heat sources: a point source $P$ at a node adds $P$ to that node's entry; a uniform volumetric source on an element adds half its total to each of the element's two nodes. Then fixed temperatures are applied and the system solved.

Two elements with conductances $g_1 = 20$ and $g_2 = 30$ mW/K connect a left end at 40 °C, a middle node and a right end at 30 °C, with a 100 mW source at the middle. The middle row is $(g_1 + g_2)T_m - g_1T_L - g_2T_R = P$, so
$$T_m = \frac{800 + 900 + 100}{50} = 36\ \text{°C}.$$
In 1D, linear FEM and finite differences coincide on uniform meshes. FEM's advantages show in 3D: meshes that follow material boundaries, local refinement and higher-order elements, at the cost of mesh generation and less regular matrices.

## Going deeper

Hughes, "The Finite Element Method: Linear Static and Dynamic Finite Element Analysis" (Dover, 2000), chapter 1; Strang and Fix, "An Analysis of the Finite Element Method" (2nd ed., Wellesley-Cambridge, 2008).
