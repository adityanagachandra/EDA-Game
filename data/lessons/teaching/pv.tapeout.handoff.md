# Tape-out

Tape-out is the handoff of the final layout to the foundry for mask making. Everything physical verification does converges here: the signed-off stream, the agreed layer map, documented waivers, and proof that nothing changed after signoff. Mistakes found after tape-out cost masks and months, so the handoff is a checklist backed by automated checks.

## What is delivered

The deliverable is a final stream, GDSII or OASIS, with a single top cell and all layers mapped by the layer map agreed with the foundry, since the stream identifies layers only by numbers. The die includes a seal ring, a continuous wall of metal and vias at the edge that stops cracks and moisture from the saw lane, and often alignment and identification marks. Supporting data includes the waivers and the signoff reports. Orphan top cells and layer-map mismatches are classic intake failures.

## The signoff checks

Physical signoff requires DRC with the full deck, LVS, antenna, density and ERC, all clean or formally waived. A waiver names the rule and the location it excuses and is agreed with the foundry; it never disables a rule globally. A final XOR between the delivered stream and the signed-off database proves that merging, renaming and conversion did not change geometry. The foundry later derives mask layers from the drawn layers by Booleans and sizing, then applies OPC and fracturing, so drawn-layer correctness is what the designer can guarantee.

## After the handoff

The foundry runs intake checks and its own DRC before committing resources and can return the design. Masks are among the most expensive items in a project, so a respin is costly. Metal-only ECOs exploit spare cells to fix bugs by changing only upper metal masks while reusing the front-end masks. Multi-project wafer shuttles share one mask set among many designs, which is how prototypes and open-source chips reach silicon. Each of these depends on verification being complete and reproducible at handoff.
