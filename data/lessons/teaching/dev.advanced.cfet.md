# CFET and the stacking outlook

A complementary FET stacks the NMOS and PMOS of a CMOS pair vertically instead of placing them side by side. Built from gate-all-around nanosheets, it removes the lateral n-to-p separation that sets standard-cell height, promising a large density step. It is still an outlook rather than a production technology, and its implications for layout, extraction, routing and modeling are what an EDA tool designer should anticipate.

## What a CFET stacks

In planar, FinFET and nanosheet CMOS, the n-type and p-type devices sit in separate regions of the cell, separated by spacing that sets part of the cell height. A CFET places a nanosheet stack of one polarity directly on top of a stack of the other, sharing one footprint. A common gate can wrap both, and the output node connects the top and bottom drains through a vertical contact.

Two fabrication approaches are discussed. Monolithic integration builds both tiers in one wafer flow, patterning the stacked layers together. Sequential integration builds the bottom tier, bonds a second silicon layer on top and then builds the upper tier, which eases alignment and lets each tier be optimized but limits the thermal budget of the upper-tier steps.

CFET is not new channel physics: each tier is a gate-all-around device. Its novelty is three-dimensional arrangement of the CMOS pair, which pushes complexity into contacts, interconnect and process integration.

## Layout, routing and extraction

Removing the n-to-p separation shrinks cell height, which reduces the number of horizontal routing tracks over each cell. Pin access gets harder, so placement and detailed routing must treat pin access as a first-class constraint.

Power delivery competes for the same scarce front-side tracks. Backside power delivery, where supply rails are routed on the back of the thinned wafer and connected to the devices from below, frees front-side tracks and gives the bottom tier a short supply path, so CFET and backside power are often discussed together.

Extraction must represent tier-aware geometry: which tier each contact and gate segment is on, vertical contacts between tiers, and capacitive coupling between stacked devices and their local interconnect, within nanometres vertically. LVS and DRC decks need tier-aware layers. Cell layout generators and place-and-route data models that assume a planar device layer need extensions for which tier a pin lives on.

## Outlook and modeling

For compact modeling, each tier is a nanosheet device, so gate-all-around channel models largely carry over. The new work is in tier-to-tier electrostatic and capacitive coupling, separate parameters for each tier (the tiers may see different thermal budgets and so different thresholds and mobility), and the parasitic network of vertical contacts.

Thermal behavior is a concern: stacking raises power per footprint and buries one tier farther from the heat-removal path, so self-heating terms and thermal-aware analysis matter more.

Density claims should be read with care. The gain from removing n-to-p spacing is real, but the net cell size depends on contact schemes, vertical interconnect and track count, and production results are not yet established. For a tool designer the safe assumption is that CFET will stress the same places FinFETs and nanosheets did, namely parasitics, pin access and variability, and add genuinely three-dimensional geometry to tools that were built around a single device layer.

## Going deeper

Ryckaert and others, "The complementary FET (CFET) for CMOS scaling beyond N3", Symposium on VLSI Technology (2018); the IEEE International Roadmap for Devices and Systems (IRDS), More Moore chapter.
