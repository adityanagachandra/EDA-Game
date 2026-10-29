# Vendor FPGA flows

Each FPGA vendor ships its own implementation suite: AMD Vivado (with Vitis for software and HLS), Intel Quartus Prime, Lattice Radiant (and Diamond for older families) and Microchip Libero SoC. They run the same stages as the open flows with signoff-grade timing data and full device and IP coverage. A tool designer should know their stages, constraint formats and artifacts, because they are the baselines open tools are measured against.

## Which tool for which device

Vivado implements AMD (Xilinx) 7-series, UltraScale and later devices; Vitis sits beside it for embedded software, HLS kernels and acceleration platforms, whose hardware is still implemented through Vivado. Quartus Prime handles Intel (Altera) families such as Cyclone and Stratix. Lattice Radiant targets Lattice's newer families, while Diamond covers older ones such as ECP5. Microchip Libero SoC handles PolarFire, SmartFusion2 and IGLOO2, with third-party synthesis built in. The tools differ in names and file formats but agree on the flow: synthesis, implementation (placement and routing), timing analysis and bitstream generation, driven by SDC-style constraints.

## Stages and constraints

In Vivado's non-project Tcl flow the steps after synthesis are opt_design, place_design, phys_opt_design, route_design and write_bitstream; constraints are XDC, which uses SDC syntax for clocks, I/O delays and timing exceptions. Quartus runs Analysis and Synthesis, then the Fitter (placement and routing), the Assembler (programming files) and the Timing Analyzer, with SDC constraints. Each step reads and writes an in-memory design database. Physical optimization after placement, such as replicating high-fanout drivers or retiming registers, uses placement information that synthesis did not have, which is a recurring pattern: later stages repair what earlier stages could only estimate.

## Checkpoints and comparisons

Vivado writes design checkpoints that store the netlist, constraints and physical state after each step, so an engineer can reopen a routed design and try a timing fix in minutes instead of rerunning hours of implementation. When comparing a vendor flow with an open one, the two use different heuristics and different timing data, so their reported $F_{max}$ values are not directly comparable; timing both routed results with one trusted timer separates algorithm from model. Vendor flows remain the choice for production designs because they cover every primitive and hard IP block with validated timing; open flows are the choice when the experiment changes the algorithm itself.

## Going deeper

AMD, "Vivado Design Suite User Guide: Implementation" (UG904) and "Using Constraints" (UG903); Intel, "Quartus Prime Pro Edition User Guide: Design Compilation"; Microchip, "Libero SoC Design Flow User Guide".
