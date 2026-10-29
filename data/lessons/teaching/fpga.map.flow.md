# FPGA synthesis flows

FPGA synthesis turns RTL into a netlist of the device's primitives: LUTs, flip-flops, carry cells, block RAMs and DSP blocks. It reuses the logic-synthesis machinery (elaboration, AIG optimization, technology mapping) but maps onto $k$-input LUTs instead of a gate library. In the open world this is a Yosys script per family, with ABC doing the LUT mapping.

## From RTL to device primitives

A Yosys FPGA script such as synth_ice40, synth_ecp5 or synth_xilinx elaborates the RTL, converts processes into flip-flops and multiplexers, flattens the hierarchy, infers memories, arithmetic and multipliers, optimizes the remaining logic, and maps it. The output is a netlist of the family's primitives, for example SB_LUT4, SB_CARRY and SB_DFF cells for iCE40, exported as JSON for nextpnr. Flattening matters: module boundaries are cut points that the LUT mapper cannot cross, so a flat netlist lets logic from different modules share LUTs. The netlist contains no placement and no routing. Those come next, and the netlist's quality, in LUT count and logic depth, bounds what they can achieve.

## LUT mapping with ABC

Yosys hands combinational logic to ABC, which builds an and-inverter graph, optimizes it and runs a cut-based LUT mapper. The mapper enumerates $k$-feasible cuts of each node, chooses a cover that minimizes LUT depth, and then recovers area on non-critical nodes with area flow and exact area. Each chosen cut becomes a $k$-input LUT whose truth table ABC computes. Telling the mapper the wrong $k$ wastes the hardware: mapping a 6-LUT device for $k = 4$ uses more LUTs and more levels. The abc9 mode in Yosys passes timing models of boxes such as carry chains and flip-flops into ABC, so mapping sees delay through the hard cells instead of treating them as opaque cut points.

## What the netlist does not know

Synthesis optimizes proxies. Unit-delay LUT depth ignores routing, which on an FPGA can exceed LUT delay, so a shallow netlist can still miss timing after placement. Reported LUT counts are not packed counts: fracturable LUTs let pairs of small functions share one element. Attributes that keep internal nets for debugging become cut boundaries and inflate LUT count. Mapping is heuristic, so a new tool version can change results on the same RTL. A tool designer should therefore judge synthesis by downstream metrics (packed area, routed delay) and should feed placement-aware information back where possible, which is exactly what timing-driven and physically aware synthesis flows try to do.

## Going deeper

Wolf, "Yosys Open SYNthesis Suite" documentation, sections on the synth_* scripts and abc9; Mishchenko, Chatterjee and Brayton, "Improvements to Technology Mapping for LUT-Based FPGAs", IEEE TCAD (2007).
