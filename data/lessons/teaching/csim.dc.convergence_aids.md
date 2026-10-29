# DC convergence aids

When plain Newton fails on an operating point, simulators fall back on a toolbox of convergence aids: junction limiting, damping, gmin stepping, source stepping, pseudo-transient continuation and user-supplied nodesets. Each aid attacks a different cause of failure and each has side effects. A tool designer chooses the order in which a simulator tries them; a user, or an agent, must diagnose which one a failing deck needs.

## Why Newton fails on circuits

Newton fails for identifiable reasons. A node with no DC path to ground, touching only capacitors or gates, has an empty row, so the matrix is singular. An exponential device evaluated at a large overshoot overflows or produces a useless Jacobian. Far from the solution, full Newton steps can cycle between two iterates. Circuits with positive feedback, such as latches, bandgaps and Schmitt triggers, have several DC solutions, and Newton may land on an unintended one or wander between their basins. A singular matrix points to topology; overflow to missing limiting; cycling to damping; an unexpected solution to the initial guess.

## Choosing an aid

Gmin stepping attaches a large conductance to ground at every node and shrinks it, making the start nearly linear. Source stepping ramps every independent source from zero, where the all-zero solution is exact, up to its full value. Pseudo-transient continuation adds artificial capacitance and integrates in time until the circuit settles, following a physical trajectory. A .nodeset gives Newton an initial guess at chosen nodes and is released after the first iterations, steering it toward the intended one of several solutions. For a bandgap that keeps landing in its zero-current state, a nodeset on the output near 1.2 V is the cleanest fix.

## What each aid costs and hides

Every aid has a price. Intermediate continuation solutions belong to a modified circuit and mean nothing until the final solve converges. Source stepping fails at a turning point, where the solution curve of a hysteretic circuit folds back, and needs arc-length continuation or pseudo-transient instead. A large global gmin set as an option stays in the circuit and corrupts high-impedance nodes. Aggressive damping turns quadratic convergence into slow linear creep: the answer is as accurate, but it costs many more factorizations. The safest fix changes the path to the solution, not the equations being solved.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter 4; the Xyce Users' Guide, chapter on homotopy and convergence; the ngspice manual, section on convergence options.
