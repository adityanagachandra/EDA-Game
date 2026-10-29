# Sources of device variability

No two transistors are identical. Dopant atoms land at random, gate edges are rough, metal grains differ, and process conditions drift across a wafer. Where variability comes from decides how it must be modeled: as a global shift shared by a whole die, as independent per-device randomness, or as a deterministic function of layout. Statistical timing, Monte Carlo simulation and SRAM yield analysis each depend on getting that classification right.

## Atomistic sources of mismatch

Three physical mechanisms dominate random device-to-device variation in modern transistors.

Random dopant fluctuation (RDF): a small doped channel contains only tens to hundreds of dopant atoms. Their number follows Poisson statistics and their positions are random, so the depletion charge and threshold vary. RDF grows as devices shrink and as channel doping rises.

Line-edge roughness (LER): lithography and etch leave gate edges that wander by a nanometre or so. The local gate length varies along the width, and through short-channel effects that becomes threshold and current variation. Fin-edge roughness plays the same role for fin width.

Metal gate granularity (MGG): the work-function metal forms grains of different crystal orientations, each with its own work function. The effective work function over a small channel depends on which grains sit above it.

Undoped thin-body devices (FinFETs, nanosheets, FD-SOI) largely suppressed RDF, leaving geometry and work-function variation as the main random terms.

## Global, local and systematic variation

Variation is classified by how it correlates.

Global (die-to-die, wafer-to-wafer, lot-to-lot) variation shifts every device on a die together: a slightly thicker oxide or different anneal for the whole wafer. It preserves matching between neighbors and is what process corners describe.

Local (within-die, random) variation is independent from device to device, even for identical neighbors. RDF, LER and MGG are local. It is what Pelgrom mismatch and LVF sigmas describe.

Systematic variation is deterministic given the layout or position: lithographic proximity effects, stress from nearby wells and shallow-trench isolation, pattern-density effects of etch and polish, and smooth gradients across a wafer. Layout-dependent effects (LDEs) are modeled with extra instance parameters such as distance to well edge or diffusion length, which extraction computes from the layout.

A Monte Carlo SPICE run must sample global parameters once per die and local parameters per device instance; mixing these up either destroys or hides matching.

## Why classification matters downstream

Different classes propagate differently through a design.

Along a path of $N$ gates, independent local deviations add in quadrature, so the path's sigma grows as $\sqrt N$ while its mean grows as $N$. The relative variation falls as $1/\sqrt N$. Global variation is fully correlated and adds linearly. Statistical timing methods (AOCV, POCV with LVF) build these rules in, and treating correlated variation as independent understates risk.

Circuits that depend on matching are dominated by local variation in the far tail. A 6T SRAM cell's read stability depends on the strength ratio of transistors inside one cell, and a large array needs that ratio to hold out to around 5 or 6 sigma. Sense amplifiers and comparators have offset set by input-pair mismatch.

Systematic effects are the most fixable: once a layout-dependent effect is modeled, it can be predicted, compensated in design, or removed by layout rules. That is why extraction tools carry LDE parameters into the netlist.

## Going deeper

Asenov, "Random dopant induced threshold voltage lowering and fluctuations in sub-0.1 µm MOSFETs", IEEE Transactions on Electron Devices (1998); Kuhn and others, "Process technology variation", IEEE Transactions on Electron Devices (2011).
