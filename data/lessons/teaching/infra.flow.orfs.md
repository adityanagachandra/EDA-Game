# OpenROAD-flow-scripts

OpenROAD-flow-scripts (ORFS) is the reference flow around OpenROAD: Makefiles and scripts that take a design from Verilog to GDS on several platforms. It runs six stages, synth, floorplan, place, cts, route and finish, each a set of steps whose order is forced by what each step needs from the ones before. A tool builder uses ORFS to reproduce results, to locate where a QoR change came from, and as a harness for new engines.

## Stages and what belongs to each

The synth stage runs Yosys and ABC to produce a gate-level netlist. The floorplan stage initializes the die and core from a utilization or explicit area, places macros, inserts tapcells and endcaps, and generates the power distribution network. The place stage runs global placement, places the IO pins, places again with the pins fixed, repairs the design by resizing and buffering, and runs detailed placement. The cts stage builds the clock tree and then repairs timing now that clock latencies are real, legalizing the cells it moved. The route stage runs global routing and then detailed routing, and inserts filler cells. The finish stage performs metal fill if the platform asks for it, writes the final reports, and merges the cell layouts into the final GDS with KLayout. Each stage saves its database under results, so a run can resume or be inspected stage by stage.

## Ordering the steps

To order a set of steps, first place each in its stage, then order steps inside a stage by data dependence. Example: the steps detailed routing, clock tree synthesis, PDN generation, repair design and global placement. By stage: PDN generation is floorplan; global placement and repair design are place; clock tree synthesis is cts; detailed routing is route. Inside place, repair design follows global placement because buffering needs cell locations to estimate wire loads. So the order is PDN generation, global placement, repair design, clock tree synthesis, detailed routing. A typical wrong option swaps two adjacent steps, for example clock tree synthesis before repair design, which looks plausible but would build a tree on a netlist that is about to change.

## Why the order is forced

The order is not convention; each step consumes an artifact of the previous one. Placement needs the rows and the power grid from the floorplan, because rails define where cells can sit. Repair design needs placed cells to estimate wire capacitance. Clock tree synthesis needs final flip-flop locations, otherwise the tree's balance is computed for the wrong sinks. Timing repair after CTS needs real clock latencies, which only exist once the tree is built. Global routing needs the clock tree's buffers and wires in place, and detailed routing needs global routing's guides. The finish stage needs every wire to compute fill and to write a complete GDS. So the step immediately after clock tree synthesis among a drawn set is the next one whose input it produces, here post-CTS timing repair if it was drawn, otherwise the earliest routing step in the set.

## Going deeper

The OpenROAD-flow-scripts documentation and its Makefile and stage scripts; the OpenROAD tool documentation.
