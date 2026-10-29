# The CMOS process flow

A chip is built by repeating a small set of unit operations several hundred times on a silicon wafer: deposit a film, pattern it with lithography, etch or implant through the pattern, clean, and sometimes polish. The order of these modules is the process flow. A tool designer needs the flow because every layer in a layout is an instruction to one of these modules, and every design rule exists because some step in that sequence cannot do better.

## Modules repeated in a fixed order

A patterning module is lithography to define a shape in photoresist, then an etch, implant or deposition that transfers it, then resist strip and clean. A leading-edge logic flow repeats such modules for dozens of mask layers. They run in three phases. The front end of line (FEOL) builds transistors: isolation trenches, wells, gate stack, spacers, source and drain. The middle of line (MOL) builds the local contacts to source, drain and gate. The back end of line (BEOL) builds the stack of copper wiring levels and vias. The order is forced by temperature: dopant activation needs about 1000 C, while finished metal and low-k films survive only about 400 C, so devices come first and wiring last.

## From drawn layers to process steps

The layout a designer hands over is a set of planar masks. The active (diffusion) layer defines where shallow trench isolation is not etched, and therefore where transistors may exist. The poly or gate layer defines gate lines; a gate-cut layer later severs them. Implant layers select which regions receive well, threshold or source-drain doping, which is why each threshold-voltage flavour is a separate layer with its own spacing rules. Contact and via layers define holes, and metal layers define trenches to be filled with copper. Each drawn layer therefore inherits the limits of its module: the resolution of its lithography, the selectivity of its etch, the conformality of its fill.

## What the flow hides from the layout

Many steps have no drawn layer: spacer deposition, epitaxial source-drain growth, silicide, dummy-gate replacement in gate-last flows, chemical-mechanical polishing, selective depositions. They still shape the final geometry and electrical behaviour. A via-to-gate short depends on film thicknesses and etch depths that the 2D layout cannot show. These undrawn steps reach design tools indirectly, through design rules, technology files and parasitic extraction models. For this reason a rule that looks arbitrary in two dimensions often has a three-dimensional cause, and a tool designer who knows the flow can tell which rules are fundamental and which might be relaxed by a DTCO study.

## Going deeper

Plummer, Deal and Griffin, "Silicon VLSI Technology: Fundamentals, Practice and Modeling", chapters 1 and 2; Sze and Lee, "Semiconductor Devices: Physics and Technology", the chapters on processing.
