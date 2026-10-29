# The views of a cell library

A cell library is not one file. It is a set of views of the same cells, each an abstraction tailored to one class of tool: the layout for manufacturing, a physical abstract for placement and routing, a timing and power model for analysis, a netlist for verification, and a behavioural model for simulation. Knowing which view each tool reads, and what each view drops, is how a tool designer diagnoses most library problems.

## The views and their readers

The full layout, in GDSII or OASIS, holds every polygon on every layer; it is used for streamout, design-rule checking and layout-versus-schematic. The LEF abstract gives each cell's size, site, symmetry, pin shapes and obstructions, plus antenna data; placers and routers read it, along with the technology LEF that defines layers, vias and routing rules. The Liberty file gives timing arcs, lookup tables, pin capacitances, power and logic functions; synthesis, static timing and power analysis read it. A transistor netlist, in CDL or SPICE format, is the reference for LVS and the input to characterization. A Verilog model gives each cell's behaviour for logic simulation.

## What each abstraction drops

Each view keeps only what its tools need. LEF hides the cell interior, keeping pins where the router may connect and obstructions where it may not route, which keeps physical data small and protects the library's layout IP. Liberty hides transistors entirely, reducing each cell to arcs and tables, so a timer can evaluate a cell by lookup instead of simulation. The Verilog model hides electrical behaviour, keeping logic and sometimes timing checks. The site definition in the technology LEF sets the unit slot that cell widths and rows are built from. Because each abstraction is lossy in a different way, no single view can answer every question, and the tools that use them never see each other's views.

## Consistency across views

A library's correctness is mostly about agreement between views. A pin capacitance in Liberty that disagrees with the layout makes every net driving that pin mistimed. A LEF pin smaller than the real metal lets the router connect badly and fail only at final DRC. A Liberty function that disagrees with the transistor netlist makes synthesis and equivalence checking wrong while LVS passes. Every one of these bugs surfaces far from its cause, in a tool that cannot see the view that is wrong. Library teams therefore generate views from a single source where possible and run cross-view checks: Liberty against SPICE, LEF against GDS, functions against netlists. A tool designer debugging odd results should suspect view mismatch early.

## Going deeper

The LEF/DEF Language Reference; the Liberty Reference Manual; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 1.
