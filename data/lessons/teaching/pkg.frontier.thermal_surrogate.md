# Fast thermal surrogates versus accurate solvers

An agent exploring package and 3D floorplans cannot run a full thermal solve for every candidate. It uses a fast engine, a compact resistance model or a learned surrogate, inside the loop, and recurses to the accurate finite-element or finite-volume solver when the fast answer is not trustworthy. Knowing what the fast engine gives up, and when to recurse, is what keeps the loop both quick and correct.

## What a fast thermal engine trades away

A compact thermal model lumps each block or die into a few nodes and solves a small network in microseconds. A learned surrogate maps a power map and stack description to a temperature map in milliseconds. Both trade resolution and generality for speed. A compact model misses hot spots smaller than its blocks and assumes fixed boundary conditions. A surrogate is accurate only near its training distribution of stacks, materials and power maps. Neither captures nonlinearity well: temperature-dependent conductivity and leakage feedback are often linearized. For ranking candidate floorplans these errors are tolerable when they shift all candidates similarly. For signoff they are not, because the limit is the absolute peak temperature.

## When to recurse to the accurate solver

The agent should call the accurate solver when the fast answer is near a decision boundary, when the input leaves the surrogate's validated range, or when candidates are too close to rank. Concretely: a predicted peak within a margin of the junction limit, a new material or cooling solution not in the training set, a hot block smaller than the compact model's resolution, or two floorplans whose predicted peaks differ by less than the surrogate's error. It should also recurse periodically to measure the surrogate's current error on real candidates. Recursion results become new training or calibration data, so the fast engine improves where the search actually spends its time.

## Measuring surrogate error honestly

A surrogate's error must be measured on held-out designs that resemble the ones the agent will explore, and reported on the quantity that matters, the peak temperature and its location, not an average over all grid points. A small mean error can hide a large error at the hot spot. The error also depends on where the search goes: an optimizer drives toward unusual designs where the surrogate may be wrong in its favour, the classic proxy-exploitation problem. The guard is to verify the final candidates with the accurate solver, track the fast-versus-accurate gap over the run, and refuse to report a surrogate-only result as a signoff temperature.

## Going deeper

Huang and others, "HotSpot: a compact thermal modeling methodology", IEEE Transactions on VLSI Systems (2006); the agentic domain concept on fast and accurate engines.
