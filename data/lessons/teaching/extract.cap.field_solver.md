# Field-solver capacitance extraction

A field solver computes capacitance from first principles: it solves the electrostatic potential around a set of conductors and integrates the resulting surface charges. It is the golden reference against which every fast extractor is calibrated, and the engine used directly on the few nets where accuracy matters most. Knowing what it computes, how its output is read, and what makes its results trustworthy is the foundation for every other extraction method.

## What the solver computes

Between conductors there is no free charge, so the potential $\phi$ satisfies Laplace's equation, $\nabla \cdot (\varepsilon \nabla \phi) = 0$, with $\phi$ fixed on each conductor's surface. Once $\phi$ is known, the surface charge density is $\sigma = -\varepsilon\, \partial \phi / \partial n$ on each conductor, and integrating it gives the conductor's charge. Because charge is linear in the boundary voltages, $Q = C V$. To get the full capacitance matrix of $N$ conductors, the solver runs $N$ solves, each with one conductor at 1 V and the rest at 0 V, and records all $N$ charges as one column. At on-chip sizes the structures are tiny compared with a wavelength, so this quasi-static electric problem is the whole capacitance story.

## Reading a solve

A solve with conductor 1 at 1 V reports charges in fC per volt, which are capacitances in fF. Suppose conductor 1 carries 8 fC and conductors 2 and 3 carry -0.75 and -2.75 fC. The driven conductor's charge is its total capacitance, $C_{11} = 8$ fF. The negative charges on the others are the negated couplings: 0.75 fF to conductor 2 and 2.75 fF to conductor 3. Charge conservation puts the remaining $8 - 0.75 - 2.75 = 4.5$ fC of negative charge on the ground plane, so the ground capacitance of conductor 1 is 4.5 fF. Field-solver post-processing does exactly this bookkeeping, and the coupling values it produces must stay attached to their aggressor nets for signal-integrity analysis.

## When the reference is golden

A solver is golden only if its inputs and discretization are. Its process description, the full dielectric stack including thin etch-stop and barrier layers, must match the foundry's, or every result errs in the same direction. Its mesh must be refined until results stop changing: a 4% shift on one refinement and 0.3% on the next means the coarse mesh was not converged. Its output must pass structural checks, such as symmetry of the matrix and non-positive off-diagonals. Because a 3D solve is costly, flows use it on critical analog or clock nets and to build and verify the tables used by fast extractors, and they approximate most of the chip by 2.5D cross-section tables where wiring is locally uniform.

## Going deeper

Nabors and White, "FastCap" (IEEE Transactions on CAD, 1991); Le Coz and Iverson, "A stochastic algorithm for high speed capacitance extraction in integrated circuits", Solid-State Electronics (1992); the method-of-moments chapter of Harrington, "Field Computation by Moment Methods" (1968).
