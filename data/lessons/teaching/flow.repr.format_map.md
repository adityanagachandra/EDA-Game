# The map of design representations and formats

A design passes through a dozen representations on its way from RTL to masks, and each one is stored in a file format that tools exchange. Each format keeps some information, drops some, and lets a particular tool compute something it could not compute otherwise. A tool designer reads formats the way a compiler writer reads intermediate representations: as contracts that say what the next stage may assume.

## The main formats and what they keep

RTL in Verilog, SystemVerilog or VHDL keeps cycle-level behaviour and drops gates and geometry. A gate-level netlist (usually structural Verilog) keeps cells and their connections and drops behaviour inside cells and all geometry. Liberty keeps each cell's function, timing and power as tables and drops its transistors. SDC keeps the designer's timing intent: clocks, input and output delays, exceptions. LEF keeps each cell's outline, pin shapes and routing blockages, and the technology's layers and rules, while dropping the cell's internals. DEF keeps a design's placement and routing. SPEF keeps extracted wire resistance and capacitance. GDSII and OASIS keep the final polygons on every mask layer and drop connectivity. A SPICE netlist keeps transistors and passive devices for circuit simulation.

## Who writes and who reads each format

Formats map onto stages. The designer writes RTL and SDC. Library characterization writes Liberty; the library provider writes LEF. Logic synthesis writes the gate netlist, which placement, routing, equivalence checking and gate-level simulation read. Floorplanning, placement and routing write DEF. Extraction writes SPEF, which static timing and power analysis read. Stream-out after routing writes GDSII or OASIS, which physical verification and mask preparation read. Delay calculation can write SDF (standard delay format), which gate-level timing simulation reads. Knowing the producer and consumer of each file tells you where a corrupt file came from and which stage it will break.

## Choosing the format for a computation

Each analysis needs a specific set of formats. Static timing analysis needs the netlist (structure), Liberty (cell delays), SDC (constraints) and, at signoff, SPEF (wire parasitics). Placement needs the netlist, LEF (cell sizes and pins) and the floorplan in DEF. Design-rule checking needs only GDSII or OASIS and the foundry's rule deck, because it is a purely geometric check. Layout-versus-schematic needs the layout and a reference netlist. A request such as "compute coupling-aware delay" therefore implies SPEF with coupling capacitances, not just a netlist. Power analysis similarly needs Liberty power tables plus switching activity from simulation, recorded as VCD or SAIF.

## Going deeper

The Liberty reference manual, the LEF/DEF language reference (distributed through Si2), IEEE 1481 for SPEF, and the OpenROAD documentation on its database (OpenDB).
