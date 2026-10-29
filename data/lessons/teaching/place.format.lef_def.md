# LEF and DEF

A placer never sees transistors. It reads LEF, which abstracts the technology and each cell into the geometry a placer and router need, and DEF, which records the design: the die, rows, instances and their locations, pins, nets and special wiring. Knowing exactly what each format keeps and drops tells a tool designer what a placement engine can and cannot compute from its inputs.

## What LEF abstracts

LEF (Library Exchange Format) comes in two parts. Technology LEF describes the process for routing: database units, each routing layer with its direction, pitch, width and spacing rules, cut layers and via definitions, and the placement SITE with its width and height. Cell LEF has one MACRO per cell: its class (core, block, pad), size, site, symmetry, every pin with its direction, use (signal, power, ground) and shapes on each layer, and OBS obstructions that block routing over internal metal. What LEF omits matters as much: no transistors, no internal wiring beyond obstructions, no timing, no power. Timing and pin capacitance come from Liberty; full layout comes from GDS. The abstract is small, which is why a placer can load a library of thousands of cells instantly.

## What DEF records

DEF (Design Exchange Format) records one design's physical state. DIEAREA gives the outline. ROW statements give each placement row's site, origin, orientation and site count; TRACKS give routing grids. COMPONENTS lists every instance with its master cell and, once placed, a status (PLACED, FIXED, COVER) with location and orientation. PINS lists the top-level I/O with layers and positions. NETS lists connectivity, and after routing the wire segments; SPECIALNETS holds power and ground rings, stripes and rails. BLOCKAGES and REGIONS restrict placement and routing. Coordinates are integers in database units, declared by UNITS DISTANCE MICRONS, so 2000 means 1 um is 2000 units. A placer reads DEF with COMPONENTS unplaced or roughly placed and writes it back with every standard cell PLACED.

## What a tool can compute from them

From LEF and DEF alone a placer can compute wirelength (from pin positions), density (from cell sizes and rows), legality (rows, sites, orientations, blockages) and a congestion estimate (from track counts and net bounding boxes). It cannot compute timing without Liberty and constraints, and it cannot check manufacturability without full layout and rules. Consistency between formats is an assumption the flow trusts: if a cell's LEF pin shape disagrees with its GDS, the router connects to the wrong place and the error surfaces only at layout-versus-schematic checking. FIXED instances are obstacles the placer must respect, and status flags let a flow freeze macros while re-placing cells. Agents operating a flow should treat these files as the contract between engines.

## Going deeper

The LEF/DEF Language Reference (Si2 and Cadence, version 5.8); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), appendix on file formats.
