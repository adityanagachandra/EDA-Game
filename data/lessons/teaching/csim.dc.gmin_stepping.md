# Gmin stepping

Gmin stepping is SPICE's main continuation method for hard operating points. It attaches a conductance gmin from every node to ground, large at first so the circuit is nearly linear and well conditioned, solves, then shrinks gmin step by step, using each solution as the starting guess for the next. The final solve uses a tiny gmin that remains in the circuit, which is worth knowing when a node has a very high impedance.

## A conductance to ground at every node

Adding gmin from every node to ground adds gmin to every diagonal of the nodal matrix. A large gmin makes the matrix strongly diagonally dominant, pins every node near 0 V, and tames the nonlinearity, so Newton converges easily. The intermediate solution belongs to a modified circuit. For a node fed by 1 mA through a 10 k$\Omega$ resistor to ground, gmin of 0.1 mS gives $v = 1/(0.1 + 0.1) = 5$ V instead of the true 10 V. That error is acceptable because the intermediate solution is only a starting point. As gmin shrinks, the solution slides continuously toward the real operating point.

## The stepping schedule

A schedule starts at some gmin, divides it by a factor after each converged solve, and ends at the final value. With a start of $10^{-3}$ S, a factor of 10 and a final value of $10^{-12}$ S, the exponents run from -3 to -12: 9 reductions and 10 solves. If a solve fails, practical implementations back off and take a smaller reduction. Each solve is a full Newton operating point, so a long schedule is expensive, and the method is used only after plain Newton fails. The schedule's parameters are options in most simulators, so an agent or user can tune them without touching code.

## The error left by the final gmin

SPICE keeps a small gmin, $10^{-12}$ S by default, across junctions, and at the end of gmin stepping, at nodes. On ordinary nodes it is negligible, but on a high-impedance node it is a real resistor of 1 T$\Omega$ to ground. A 1 V source driving a node through 1 G$\Omega$ with only gmin to ground forms a divider: $v = 1/(1 + 10^{9} \times 10^{-12}) = 0.999$ V, an error of about 1 mV. Leakage paths, sample-and-hold nodes and floating gates sit at exactly these impedances, so their simulated voltages may depend on gmin. Raising gmin to help convergence makes this error larger.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter 4; the ngspice manual, sections on gmin and convergence aids.
