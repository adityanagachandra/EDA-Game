# Open-source EDA tools

Open-source EDA has grown from academic prototypes into a complete RTL-to-GDSII flow that has taped out real chips. For a tool designer it is a working laboratory: every engine's source is readable, every intermediate file is inspectable, and a change can be tested from RTL to layout. Knowing what each open tool does, and what it reads and writes, is the practical starting point for writing EDA software.

## The digital implementation chain

Yosys reads Verilog RTL, elaborates and optimizes it, and writes a gate-level netlist; it hands logic optimization and technology mapping to ABC, which works on AIGs. OpenROAD takes that netlist with LEF, Liberty and SDC and runs floorplanning, global placement (based on RePlAce), detailed placement, clock tree synthesis (TritonCTS), global routing (FastRoute) and detailed routing (TritonRoute), storing everything in its OpenDB database and writing DEF. OpenSTA, embedded in OpenROAD and usable alone, reads the netlist, Liberty, SDC and SPEF and reports timing. OpenLane and the OpenROAD flow scripts wire these into push-button flows.

## Layout, verification and simulation

KLayout views and edits GDSII and OASIS and runs DRC and LVS scripts. Magic is a layout editor with interactive DRC and parasitic extraction. Netgen compares two netlists for LVS. ngspice and Xyce (from Sandia National Laboratories, designed for parallel simulation) run SPICE netlists. Verilator compiles Verilog into a fast cycle-based C++ model; Icarus Verilog is an event-driven simulator; cocotb writes testbenches in Python. SymbiYosys drives formal engines on SystemVerilog assertions. nextpnr and VPR place and route FPGA designs, VPR being the research platform of the VTR project. Together these tools cover simulation, formal checking, layout and LVS, so a complete chip can be verified without commercial licences.

## What open tools are good for

Open tools trail commercial ones in QoR on the hardest advanced-node designs, and they are certified for fewer processes. Their strengths are elsewhere: transparency, scriptability, and freedom to modify. Research results can be reproduced and compared; a new placer heuristic can be dropped into OpenROAD and measured on public benchmarks; a flow can run thousands of times in the cloud without licence limits, which is what machine-learning and agentic experiments need. Open PDKs at mature nodes, such as SKY130 and GF180MCU, closed the last gap by allowing real fabrication. They are also the natural base for agent-first experiments, since an agent can read and modify the engine source.

## Going deeper

Ajayi et al., "Toward an Open-Source Digital Flow: First Learnings from the OpenROAD Project" (DAC 2019); Wolf, "Yosys Open SYnthesis Suite"; the VTR project documentation.
