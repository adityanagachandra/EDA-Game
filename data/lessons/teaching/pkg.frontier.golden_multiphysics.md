# Golden multiphysics checks

Before an agent's change to a multiphysics engine or flow is accepted, it must pass golden and differential tests. Golden tests compare the engine with known answers; differential tests compare the changed engine with the baseline on reference designs, proving that only the intended results moved. For thermal, mechanical and electromigration engines, the known answers come from closed-form physics and from measurements.

## Golden cases from closed-form physics

Good golden tests for multiphysics engines have exact answers. For thermal: a 1-D slab with fixed temperatures, whose resistance is $L/(kA)$; a chain with distributed heating, whose far-end rise is $RqN(N+1)/2$; a lumped body's step response $PR(1 - e^{-t/\tau})$. For mechanics: a bilayer strip, whose curvature follows Timoshenko's formula; a constrained thin film, whose stress is $E\Delta\alpha\Delta T/(1-\nu)$. For coupling: a self-heating resistor, whose converged temperature has a closed form and whose runaway current is known. Each should pass to a tolerance tied to mesh resolution, and the suite should include mesh-refinement checks showing the expected convergence order.

## Differential checks on reference designs

A differential test runs baseline and changed engines on reference packages and compares every reported quantity: per-die peak temperatures and their locations, warpage across the temperature profile, joint strains, current densities. A change that should affect only one thing, say a new creep model for one solder, must leave thermal results identical to the last bit and move only the mechanical results of assemblies that use that solder. Any unexpected movement elsewhere is a regression or a hidden coupling. The check needs deterministic engines, seeded meshing and stable ordering, otherwise numerical noise hides real differences.

## Measured ground truth

Closed forms cover simple geometry; measurements anchor real packages. Thermal test dies with embedded heaters and sensors give junction temperatures for known power maps. Shadow moiré gives warpage across the reflow profile. Accelerated cycling gives joint lives. A golden suite should include a few such measured cases with their uncertainty, and the engine passes when its prediction falls within that uncertainty. These cases also detect drift in material libraries, a common silent failure: a changed modulus table can pass every closed-form test and still mispredict a real package. Measurements are design and process data, so they stay inside the IP boundary, with only pass or fail results shared.

## Going deeper

Oberkampf and Roy, "Verification and Validation in Scientific Computing" (2010); the agentic domain concept on golden and differential testing.
