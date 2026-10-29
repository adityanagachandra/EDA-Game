# The extractor's core and change surface

An extractor built for an agent-first flow is split into a small immutable core and a change surface that agents and engineers can adjust without touching the core. For extraction the core is the field-solver mathematics and netlist correctness; the change surface is the accuracy-versus-speed mode, net and layer selection, the pattern library and calibration. Placing each customer request in the right layer and tier, and putting the seam where requests actually vary, is what keeps core edits rare.

## Core and change surface

The four-layer model orders an engine's parts by how hard they are to change. The immutable core of an extractor holds what makes it correct: the discretized Laplace solve and its numerical kernel, the assembly of the capacitance matrix, and netlist correctness, meaning every resistor and capacitor is attached to the right net and node. These are shared by every customer and change only with human approval and full regression. Strategy covers the pattern-matching and fracturing policies. Parameters cover the accuracy mode, the net and layer selection, thresholds for coupling filtering and calibration tables. Observation hooks expose per-net parasitic summaries, couplings by aggressor and the estimated error of the mode used, so an agent can decide what to re-extract.

## Tiers for extraction requests

Requests land in the smallest safe tier. Extracting clock nets with the 3D solver and the rest with patterns is Tier 1, a selection setting. Refitting pattern coefficients to new foundry golden data is Tier 1, a calibration table change validated against the golden results. A new cell architecture whose via-rail structure matches no existing pattern needs new recognition code at the pattern-library seam: Tier 2, a plug-in. Changing the mesh-refinement algorithm inside the 3D solver is Tier 3, a core edit with a reviewed diff and regression against every golden result. A seam is placed well when the common requests, new processes, new layout styles and new accuracy targets, are all Tier 1 or Tier 2.

## Designing the capacitance-model seam

The seam between geometry and capacitance sits at a capacitance-model interface: the core fractures geometry and builds connectivity, then asks a model for per-length ground and coupling values for each cross-section. To host table, calibrated and surrogate models alike, the seam must pass the layer, width, spacing to each neighbour and that neighbour's net, the layers above and below, local density and the corner. It must be coarse, called per batch of segments with plain arrays rather than per segment through a virtual call, so the core's hot loop stays fast. It must be deterministic: the same layout must give the same parasitics on every run, or differential checks cannot separate an ECO's effect from noise.

## Going deeper

The OpenRCX extractor in OpenROAD, whose pattern-based model files illustrate a calibration data seam; Ousterhout, "A Philosophy of Software Design" (2nd ed., 2021), on deep modules and interfaces.
