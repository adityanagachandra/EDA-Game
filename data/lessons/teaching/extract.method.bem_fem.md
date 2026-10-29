# Boundary-element versus finite-element solvers

Deterministic field solvers come in two main families. Boundary-element methods discretize only conductor surfaces and produce small, dense systems. Finite-element and finite-difference methods mesh the whole volume and produce large, sparse systems. Each fits different structures, and an extraction engine designer must understand their unknown counts, memory and solve costs to pick and accelerate the right one.

## Surface panels

A boundary-element method (BEM), a form of the method of moments, writes the potential as an integral of unknown surface charge against the free-space Green's function. It discretizes each conductor surface into panels with a constant charge each, and enforces the known potential at each panel. The unknowns are the panel charges, so their count is the surface area divided by the panel area. A box-shaped wire 1.4 um long, 0.05 um wide and 0.15 um thick, cut into 0.05 um square panels, is $28 \times 1 \times 3$ in panel units, with $2(28 + 84 + 3) = 230$ panels. Only surfaces are meshed, so the unknown count grows with the square of feature size, not the cube.

## Dense matrices

Every panel's charge affects the potential at every other panel, so the BEM system matrix is dense: $N^2$ entries for $N$ panels. In double precision that is $8N^2$ bytes; 17000 panels need $8 \times 17000^2 = 2.3 \times 10^9$ bytes, 2312 MB. A direct LU solve costs $O(N^3)$ operations. Practical BEM extractors therefore never form the matrix: they solve iteratively with GMRES and compute matrix-vector products with acceleration such as the fast multipole method (as in FastCap) or precorrected FFT, bringing cost near $O(N)$ or $O(N \log N)$ per iteration. A preconditioner is needed to keep the iteration count low.

## Volume meshes and sparsity

A finite-element method (FEM) meshes the dielectric volume and solves for the potential at mesh nodes. Each node couples only to its neighbours, so the matrix is sparse: a tetrahedral or hexahedral mesh has on the order of 7 to 27 nonzeros per row. FEM uses many more unknowns than BEM for the same structure but far fewer per unknown. With 2000 BEM panels against 100000 FEM unknowns at 27 nonzeros each, the stored entries are $4 \times 10^6$ versus $2.7 \times 10^6$, a ratio of 1.48. FEM handles complex, layered dielectrics and nonuniform permittivity naturally, since material changes are local, whereas BEM needs a layered Green's function or extra unknowns on dielectric interfaces. Both, and random-walk methods, serve as golden references for pattern tables.

## Going deeper

Nabors and White, "FastCap: a multipole accelerated 3-D capacitance extraction program", IEEE Transactions on CAD (1991); Phillips and White, "A precorrected-FFT method for electrostatic analysis of complicated 3-D structures", IEEE Transactions on CAD (1997); Jin, "The Finite Element Method in Electromagnetics" (3rd ed., 2014).
