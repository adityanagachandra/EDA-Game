# Yosys and ABC in an open flow

Every open digital flow starts with Yosys, which reads Verilog, elaborates and optimizes it in its RTLIL intermediate representation, and hands the combinational logic to ABC for optimization and technology mapping against a Liberty library. A tool builder needs the script's structure, the exact handoff between the two tools, and the common failure signatures, because synthesis problems surface downstream as strange placement or timing.

## The shape of a Yosys synthesis script

A typical open-flow script reads the Verilog, then runs hierarchy with the top module named, which elaborates parameters and reports undefined modules. proc lowers always blocks into explicit multiplexers, flip-flops and latches. opt, fsm and memory passes simplify logic, extract and recode state machines, and map memories. techmap lowers generic cells into a fine-grained gate library. dfflibmap with the Liberty file maps the internal flip-flop cells onto real library flip-flops. abc with the Liberty file maps the remaining combinational logic. opt_clean removes dangling wires and cells, and write_verilog emits the structural netlist. Flows add passes for their downstream tools, such as tie-cell insertion for constants or buffer insertion for direct port-to-port assigns. Yosys also has a synth command that bundles a standard sequence; open flows usually write the steps out to control them.

## What ABC receives and returns

When Yosys runs abc, it extracts the combinational logic, writes it out for ABC, and reads the mapped result back, stitching it into the netlist between the flip-flops and other cells Yosys keeps. ABC represents logic internally as and-inverter graphs, optimizes them with rewriting, refactoring and balancing, then maps to the library using each cell's Boolean function, area and delay from the Liberty file. Geometry from LEF and clocks from SDC never enter. If a flow passes a clock period or delay target, ABC favors faster structures, often at more area; without one, it leans toward area. An example: a 32-bit adder mapped with no target tends toward a compact, slower structure; with an aggressive target ABC spends more cells to shorten the critical path, and the reported area rises.

## Diagnosing open synthesis failures

Three failures recur. First, "module not defined" for an SRAM or other hard macro: there is no RTL to synthesize, so the flow must supply a blackbox declaration or a Liberty model so Yosys keeps the instance as a black box. Second, a netlist far smaller than expected, with outputs tied to constants: logic that drives no observable output is legally removed, so an unconnected port, a typo in a port name or a missing keep attribute shows up as missing logic. Third, large area differences between two Liberty files from one PDK: mapping can only use the cells offered, so a library with richer complex gates or different sizes changes the result. A fourth, softer symptom is a netlist with plain assigns between ports that a physical tool rejects; flows insert buffers to give every connection a driving cell.

## Going deeper

The Yosys manual (passes proc, techmap, dfflibmap, abc); the ABC documentation and Mishchenko's papers on AIG rewriting and mapping; the OpenROAD-flow-scripts synthesis scripts.
