# The PCB design flow

A printed circuit board flow takes a circuit drawn from purchased parts and turns it into copper, holes and assembly instructions. It starts where the chip flow ends: the silicon is already packaged, so the board tool never synthesizes logic. It captures connectivity, binds every part to a physical footprint, places and routes a few hundred to a few thousand packages on a few to a few dozen layers, checks the result and writes manufacturing data. A tool designer needs this map to see which engines exist, what each consumes and where errors leak between them.

## Stages of the board flow

The stages run in order. **Schematic capture** records components and nets. **Footprint assignment** binds each symbol to a land pattern. **Netlist transfer** moves connectivity and footprints into layout. **Placement** positions the packages, **routing** draws traces and vias on the stackup, **DRC** checks geometry and connectivity against the fabricator's rules, and **output generation** writes the fabrication data (Gerber or ODB++, drill files, stackup drawing) and the assembly data (BOM, centroid file, paste layers). Unlike the chip flow there is no synthesis, no standard-cell library and no foundry PDK. The parts are fixed packages, and the rules come from a board fabricator's process capability, usually expressed against IPC standards. A small board might be 2 layers and 50 parts; a server board might be 20 or more layers and several thousand parts.

## What crosses each boundary

Each handoff is a representation that preserves some facts and drops others. The **netlist** carries nets as lists of reference-designator and pin-number pairs, plus the footprint name of each part, and nothing geometric. The **board database** adds outline, stackup, placed footprints, traces, vias, pours and rules. The **fabrication outputs** are images and drill programs per layer; classic Gerber has no connectivity, so a separate IPC-D-356 netlist goes to bare-board test. The **assembly outputs** split by audience: the BOM tells purchasing what to buy, the centroid file tells the placement machine where each part goes and at what rotation, and the paste layers become the stencil. Two consumers, a bare-board fabricator and an assembler, see entirely different subsets of the design.

## Iteration and failure propagation

Board flows iterate just as chip flows do, and errors propagate downstream silently. A pinout error in a footprint passes every layout check, because the copper matches the netlist exactly; it fails only on the bench. A placement that crowds a BGA against a connector shows up later as unroutable channels, so the fix is upstream at placement. A stackup change from the fabricator after layout moves every controlled impedance, because impedance rises with dielectric height, so trace widths must be recomputed. The engineering rule is to fix a problem in the representation where it originates and push it forward: a mirrored connector is fixed in the library or schematic, then forward-annotated, never patched only in copper, because the next netlist update would undo the patch.

## Going deeper

Bogatin, "Signal and Power Integrity Simplified" (3rd ed., 2018), chapter 1, for why board interconnect is an engineering problem in its own right. The KiCad documentation, "Getting Started in KiCad", for a complete open-source walk through the flow.
