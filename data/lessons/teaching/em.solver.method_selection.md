# Choosing a field solver: MoM, FEM, FDTD and PEEC

No single field solver is best for every EDA task. Boundary-element methods, finite elements, FDTD and PEEC discretize different things, work in different domains and scale differently. Choosing the right one for a capacitance table, an on-chip inductance model or a package S-parameter sweep is an engineering decision that tool architects make once, and that an agent selecting among solver skills must make every time.

## Surface versus volume discretization

The method of moments meshes only conductor surfaces and uses a Green's function for the medium. In a homogeneous or simply layered dielectric this gives far fewer unknowns than meshing the volume, but the matrix is dense and must be compressed with fast multipole or hierarchical methods to scale.

The finite-element method meshes the whole volume, assigning a material to each element, so arbitrary inhomogeneous dielectrics and complex shapes are natural; the matrix is sparse but large, and open regions must be truncated with absorbing boundaries. Finite differences sit at the simple end of volume methods. A rule of thumb: surface methods for many conductors in simple media, volume methods for complex material stacks. Layered Green's functions extend MoM to planar stacks, which is why it remains popular for package and board extraction.

## Time domain versus frequency domain

FDTD steps fields in time: one run with a broadband pulse yields the response over the whole band after a Fourier transform, but the step is bounded by the Courant limit and high-Q structures need long runs. Frequency-domain FEM and MoM solve one complex linear system per frequency; they excel at resonant structures and narrow bands, and adaptive frequency sweeps fit rational models from a few solves to cover a band cheaply.

Example decision: a connector characterized from DC to 50 GHz with many resonant features can go either way, and the choice depends on geometry detail and Q. A narrowband filter with very high Q favours frequency domain. A structure with nonlinear devices, or one that must be co-simulated with circuits in transient, favours time domain. Both must agree on a common golden structure, which is how teams validate one against the other.

## Matching the solver to the extraction task

Full-chip parasitic extraction uses pattern-matched tables calibrated by 2D and 3D field solvers, because billions of segments cannot each be solved; accurate 3D solves are reserved for critical nets. On-chip inductance and frequency-dependent resistance use magnetoquasistatic methods such as PEEC, which turns conductor filaments into partial resistances and partial inductances with mutual coupling that plug directly into circuit simulation. Package, board and connector S-parameters at tens of gigahertz use 3D full-wave FEM, FDTD or MoM, since structures there are a sizeable fraction of a wavelength.

Example: a team needs the inductance of a power grid with no radiation: PEEC. A team needs S-parameters of a via transition to 56 GHz: 3D full-wave. A team needs capacitance of every net overnight: calibrated tables. Choosing the full-wave solver for the first or third task wastes orders of magnitude of compute.

## Going deeper

Ruehli, "Equivalent circuit models for three-dimensional multiconductor systems", IEEE Transactions on Microwave Theory and Techniques, 1974; Jin, "The Finite Element Method in Electromagnetics" (3rd ed.); Kamon, Tsuk and White, "FastHenry: a multipole-accelerated 3-D inductance extraction program", 1994.
