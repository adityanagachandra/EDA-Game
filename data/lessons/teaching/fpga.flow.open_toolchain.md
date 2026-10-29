# The open FPGA toolchain

An open FPGA toolchain is a complete, source-available flow from Verilog to bitstream: Yosys for synthesis, nextpnr or VPR for packing, placement and routing, and a documented bitstream database with its packing tool. It is what lets a tool designer read, change and measure every stage, and it is the platform for most FPGA CAD research and for agent experiments on real devices.

## The stages and their files

For iCE40 the flow is Yosys with synth_ice40, writing a JSON netlist; nextpnr-ice40, reading the JSON netlist, a pin constraint file and the device database, and writing a textual configuration; and icepack, turning that into the binary bitstream. ECP5 is the same with synth_ecp5, nextpnr-ecp5 and ecppack over the Project Trellis database. The VTR flow is the research counterpart: a Verilog front end, ABC for LUT mapping, and VPR for packing, placement, routing and timing on an architecture written in XML. In VPR-based flows for real devices, the routed design is written as FASM, a list of enabled configuration features, which a database turns into frames.

## Real devices versus architecture models

nextpnr's architectures are real chips: the device database fixes every tile, wire and switch, and the output must load and run. VPR's architecture is a parameter: the XML file defines logic blocks, segment lengths, switch patterns and hard blocks, and VPR builds the routing-resource graph from it, so a researcher can evaluate an FPGA that has never been built. That is why VPR defaults to searching for the minimum channel width, an architecture metric, while nextpnr routes on the fixed fabric and reports timing from the device's delay data. Both share the same algorithms: annealing or analytical placement, and PathFinder-style negotiated routing. Choosing the tool means choosing between bit-exact devices and architectural freedom.

## Running fair experiments

Open flows make controlled experiments possible, but only if they are controlled. Placement is randomized, so a result needs several seeds, and a design near its limit can pass timing on one seed and fail on another. A change to one stage must hold every other stage fixed: the same benchmarks, architecture file, and downstream placer and router when evaluating a packer. A router change that lowers minimum channel width must also report critical-path delay and run time at a relaxed width, since routability can be bought with delay. And the bitstream database is itself data that can be wrong, so its coverage limits what an open flow can target. These are the habits that turn a toolchain into an instrument.

## Going deeper

Shah et al., "Yosys+nextpnr: An Open Source Framework from Verilog to Bitstream for Commercial FPGAs", FCCM (2019); Murray et al., "VTR 8", ACM TRETS (2020); the nextpnr and Yosys documentation.
