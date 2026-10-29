# OpenDB

OpenDB is the design database at the center of OpenROAD. It descends from the Athena database and its schema follows LEF and DEF: technology layers and vias, library masters and pins, and a block of instances, nets, terminals, rows and wires. Every OpenROAD engine, from floorplanning to detailed routing, reads and edits the same OpenDB objects in one process. Knowing its object model is the entry point to reading or extending any OpenROAD tool.

## Technology, libraries and blocks

OpenDB separates data by origin. A dbDatabase holds one dbTech, built from technology LEF, with dbTechLayer objects for each routing and cut layer (direction, pitch, width, spacing) and the via definitions. Library LEF fills dbLib objects holding dbMaster cells, each with dbMTerm pins and their shapes and obstructions, plus sites. The design is a dbChip with a top dbBlock, built from DEF or from a Verilog netlist plus floorplan commands. The block owns dbInst instances, dbNet nets, dbBTerm ports, dbRow rows and the routing. This mirrors LEF and DEF so closely that reading DEF is almost a direct mapping of sections to object creation, and writing DEF back is a traversal. Timing data is not in OpenDB: Liberty is read by OpenSTA, which keeps its own library model and links to OpenDB instances.

## Navigating instances, terminals and nets

The links follow the object model. A dbInst points to its dbMaster and owns one dbITerm per master pin. Each dbITerm points to its dbMTerm (giving the pin's direction and shapes) and to the dbNet it connects, if any. A dbNet holds its set of dbITerms and dbBTerms. So to list the sinks of net n1, iterate the net's ITerms, keep those whose master terminal is an input, and read each one's instance. To get a cell's width, ask its master, not the instance: the instance stores only location, orientation and placement status. Engines that need to react to edits register a callback object on the block and are told when instances move or terminals connect. The Tcl and Python interfaces expose the same objects, so scripts can walk the database exactly as C++ code does.

## Flat blocks, wires and checkpoints

In the usual OpenROAD flow the block is flat: a hierarchical Verilog netlist is linked into a single level of leaf instances, with hierarchy preserved in names and, in recent versions, optional module objects. Connectivity and geometry are distinct objects: a dbNet is logical connectivity that exists from synthesis onward, while a dbWire, written by the detailed router, encodes the routed shapes and vias. Power and ground live as special nets, from DEF SPECIALNETS, with special wires for stripes and rails. Finally, OpenDB has a native binary file format. Saving it between flow stages, as OpenROAD-flow-scripts does, keeps every object and property and reloads far faster than DEF, which is an interchange format with a fixed vocabulary. A stage that crashes can be rerun from the previous checkpoint.

## Going deeper

The OpenDB documentation and source in the OpenROAD repository (src/odb); Ajayi et al., "Toward an Open-Source Digital Flow: First Learnings from the OpenROAD Project", DAC 2019.
