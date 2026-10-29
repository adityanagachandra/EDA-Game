# Graph-based versus path-based analysis

A timer can analyze in two modes. Graph-based analysis keeps one worst value per pin and times the whole design in one pass. Path-based analysis takes specific paths and re-times each with the values that actually travel along it. GBA is fast and safe but pessimistic; PBA is accurate for the paths it touches but expensive. Knowing what each mode merges and when to switch is part of designing a closure flow.

## What graph-based analysis merges

At every pin GBA stores the worst arrival and the worst slew, each over all incoming arcs, computed separately. A downstream arc's delay is then evaluated at a slew that may belong to a different input than the latest arrival, so the arrival can exceed every real path. GBA also uses conservative values wherever a quantity depends on the whole path: advanced OCV derates indexed by logic depth must use a safe depth for the pin, since one pin lies on paths of many depths. Each merge can only overestimate setup delay, so GBA is a safe bound. Its cost grows with the graph, not with the number of paths, which makes it the only mode that can analyze a whole chip.

## What path-based analysis recovers

PBA takes a path, usually one that fails under GBA, and recomputes every stage with that path's own input slews, its true logic depth for AOCV, and its own common clock point for CPPR. Its slack is always at least the GBA slack for the same path. A path that shows -12 ps under GBA and $+4$ ps under PBA actually meets timing; the 16 ps was merge pessimism. Because every path that passes under GBA also passes under PBA, a flow only needs to re-time paths whose GBA slack is below a threshold. PBA is far slower per path, and a change to any cell on a path makes its PBA result stale.

## Using both in a flow

Optimizers run on GBA because its per-pin values update incrementally after each change; PBA would have to re-time many paths after every move. Signoff then applies PBA to the worst paths to avoid fixing violations that are not real, which saves area and power. Exhaustive PBA over every path would be exact, but path counts grow exponentially with reconvergence, so tools enumerate the worst paths in slack order and stop at a limit. A practical consequence for tool builders: if an optimizer sees GBA slack while signoff uses PBA slack, the optimizer over-fixes. Some flows therefore apply PBA-based adjustments during the final optimization steps.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter on path-based analysis; the OpenSTA and PrimeTime documentation on path-based analysis.
