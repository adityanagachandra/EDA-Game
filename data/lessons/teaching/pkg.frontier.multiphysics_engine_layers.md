# Multiphysics engines in the four-layer model

Package and 3D-IC analysis runs on multiphysics engines: thermal solvers, mechanical FEA, power-integrity and electromigration checkers, and the coupling between them. An agent that customizes these engines must know which part is immutable core and which is change surface. The four-layer engine model sorts every request into core, policy, parameter or hook, and the change tiers say how each is made and reviewed.

## Core and change surface of a multiphysics engine

In a thermal or mechanical solver the core is the discretization and the numerical kernel: assembling the conductance or stiffness matrix from the mesh, enforcing boundary conditions, and the sparse solver with its convergence test. Energy conservation and the material law are invariants. The change surface is everything that varies by customer and package: material libraries with temperature-dependent properties, mesh-refinement policies, which coupling scheme to use (loose iteration or monolithic Newton), convergence tolerances, the power maps imported from chip tools, and the reports to produce. Hooks expose live state such as per-die peak temperature, the coupling residual and the joint strain, so an agent can watch a run and decide what to do next.

## Placing a request in a tier

Tier 1 is configuration: a new material card, a tighter tolerance, a different heat-transfer coefficient, a new report. No code changes. Tier 2 is a plug-in at a defined seam: a new material model such as a creep law for a new solder, a new power-map importer, or a new boundary-condition type such as a two-phase cold plate, written against the engine's material or boundary interface. Tier 3 is a core edit: changing how the matrix is assembled, or replacing the solver. A request to model a new interface material with a known conductivity is Tier 1; a request to model its pressure-dependent bond line needs a Tier 2 material plug-in. A seam is well placed when such requests never reach Tier 3.

## Seam design for coupled physics

Coupling is where multiphysics engines most need good seams. A clean design exposes a field-exchange interface: each single-physics solver publishes fields such as power density, temperature and displacement on its own mesh, and a coupler maps them between meshes and decides when to iterate. Then a new coupling, for example temperature-dependent electrical resistance feeding back into power, is a coupler policy rather than an edit inside both solvers. The seam must carry enough context: mesh, units, time, and validity ranges of each material. Calls should be coarse, per field and per iteration, not per element, so the kernels stay fast. Everything must be deterministic so that a differential test can compare runs.

## Going deeper

The agent-first EDA concepts in this curriculum's agentic domain; Keyes and others, "Multiphysics simulations: challenges and opportunities", International Journal of High Performance Computing Applications (2013).
