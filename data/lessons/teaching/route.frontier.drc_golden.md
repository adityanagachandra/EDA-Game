# DRC-clean as the router's golden test

A router change is accepted only if routed reference designs come out clean. DRC-clean, connected and within a quality-of-results tolerance is the router's golden test, and a differential comparison against the baseline proves the change touched only what it was meant to. These checks are what let an agent edit routing policy at all.

## Wirelength against the baseline

Each golden design has a stored baseline routing. A new build must route it with zero DRC violations and zero opens, the hard gates, and then its quality metrics are compared with the baseline: wirelength, via count, overflow, timing. The wirelength change is $(WL_{new} - WL_{base}) / WL_{base} \times 100\%$.

Example: a baseline of 264972 um and a new result of 265025 um is a change of 0.02%. Hard gates come first: a shorter routing with one short is a failure, not an improvement. Tolerances are needed because routers are heuristic, but they must be tight enough to catch real regressions.

## Differential checks

A golden test asks whether the result is good; a differential test asks whether it changed only where intended. Route the reference design with and without the change and compare per net. Example: a plug-in meant to affect clock nets only must leave every signal net's route identical. If two signal nets changed wirelength, the check fails, even if total wirelength improved.

Unintended changes reveal one of two problems: the seam leaks behaviour outside its scope (the plug-in's cost touched edges used by other nets), or the router is not deterministic. Both must be fixed before promotion, because otherwise nobody can tell whether later differences come from the change or from noise.

## The golden suite

One design is not enough: a suite covers dense and sparse regions, macros, many layers and few, different rule decks. A change is promoted only when every design passes every gate. Example: with a 1% wirelength tolerance, results of +2.7%, -1.9%, -2.2%, +0.5% with two opens, and +0.5% clean fail four of five designs.

For an agent this suite is the reward signal, and it must be hard to game. Checking DRC with the router's own fast checker invites changes that satisfy the checker but not the signoff engine; skipping connectivity invites deleting wires to remove violations. Golden criteria therefore use the signoff DRC engine, an independent connectivity check and timing, and they are versioned together with the baselines.

## Going deeper

The OpenROAD regression test suite as an example.
