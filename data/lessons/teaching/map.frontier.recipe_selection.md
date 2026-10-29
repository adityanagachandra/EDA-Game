# Agents choosing mapping recipes

Most mapping wins in practice come not from new algorithms but from choosing the right script: which mapper, which K, whether to compute choices, how many area-recovery passes. An agent can make that choice from manifests and measured QoR, if the recipes are packaged as skills with honest descriptions and the loop is grounded in reference engines rather than proxies.

## Recipes as skills with manifests

Package each mapping recipe as a skill: the script, its build and dependencies, golden tests, and a short manifest stating the target (LUT or standard cell), objective (depth-first, delay-first with area recovery, area-only), inputs and outputs, tunable parameters (K, priority cuts per node, delay target, number of recovery passes), and known strengths and weaknesses. The agent selects by reading manifests, a cheap text decision: an FPGA target excludes standard-cell recipes, a failing timing report points to a delay-first recipe, and a hard area cap points to area-oriented ones.

## Running and comparing candidates

An agent can run several recipes in parallel and keep the best, but only among valid results: each candidate must pass equivalence checking and meet the timing constraint before area is compared, and all must run on the same design, library and constraints. Fast estimates, such as the mapper's own depth and area or a quick timer, filter recipes cheaply; the reference timer, ideally after placement, decides. The loop needs a stopping rule, an iteration budget and a minimum improvement, or it oscillates between recipes that each fix what the other broke.

## Learning which recipe to pick

Each run is data: design features (size, depth, fanout distribution, datapath versus control), recipe, parameters and measured QoR. Recording when a recipe wins, in its manifest and in the domain model, is the skill flywheel: the next selection starts from experience. Rewards for learned selection must use the real objectives measured by reference engines. A reward on cell count alone invites recipes that use few large cells and hurt area and delay; a reward on the mapper's own delay estimate invites results that collapse after placement. Tool-grounded rewards measure what signoff will measure.

## Going deeper

The ABC and Yosys documentation for the commands recipes are built from; Mishchenko et al., "Combinational and sequential mapping with priority cuts", ICCAD 2007, for the parameters a recipe exposes.
