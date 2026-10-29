# Table models and fast SPICE accuracy

Besides partitioning, fast SPICE engines save time by simplifying what they evaluate: table models in place of compact models, relaxed tolerances on digital-like blocks, reuse of identical subcircuit solutions, and reduced parasitic networks. These choices are why full-chip transistor-level simulation is possible at all, and also why its results must be trusted selectively. The designer's skill is knowing when the fast answer is good enough and when to recurse into accurate SPICE.

## Replacing compact models with tables

A BSIM-class model evaluates long sequences of equations per device per Newton iteration. A table model precomputes currents and charges on a grid of terminal voltages and interpolates, costing a few multiply-adds. Interpolation must keep derivatives continuous: piecewise-linear tables have piecewise-constant conductances that jump at grid lines, which makes Newton oscillate across them, so engines use splines or other smooth schemes. A table is valid only where it was built: other temperatures, geometries or corners need their own tables or extra dimensions. Grid resolution trades memory and build time against accuracy, and the subthreshold region, where current changes by decades over a few hundred millivolts, is the hardest to resolve.

## Relaxed tolerances and reuse

Digital partitions swing rail to rail, so millivolt errors rarely change edge timing or power, and fast engines relax Newton and step tolerances there. Memory arrays contain thousands of identical cells in identical states; isomorphism detection solves one and reuses the result for all with the same boundary conditions. Post-layout netlists carry millions of parasitic resistors and capacitors, which are reduced before simulation by node elimination or moment-matching projection, keeping delay behavior with far fewer elements. These choices are bundled into accuracy or speed modes, so a single setting moves tolerances, table resolution and partitioning aggressiveness together, which is a classic Tier 1 control for a user or an agent.

## When fast SPICE is not enough

Fast SPICE suits full-chip power, digital timing of custom blocks and functional checks of large mixed designs. It is the wrong tool where the answer depends on small differences: comparator offset to 0.1 mV, leakage of a sleeping memory, distortion of an analog front end, or the stability of a tightly coupled feedback loop that a partitioner may cut and treat as weak. The productive practice is two-level: run the fast engine broadly, then recurse into accurate SPICE on the few blocks or paths where the fast engine's error budget could change a decision. Sign-off of precision analog metrics stays on accurate simulation.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter on fast simulation; Rabbat, Sangiovanni-Vincentelli and Hsieh, "A multilevel Newton algorithm with macromodeling and latency for the analysis of large-scale nonlinear circuits in the time domain" (IEEE TCAS, 1979).
