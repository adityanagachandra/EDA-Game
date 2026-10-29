# The FPGA implementation flow

An FPGA flow starts from the same RTL as an ASIC flow but ends in a bitstream that configures a fixed chip. Synthesis maps logic to lookup tables (LUTs), flip-flops and hard blocks; packing groups them into the fabric's logic blocks; placement assigns those blocks to sites on the device; routing chooses programmable switches in a fixed wiring network; and the bitstream generator writes the configuration. The algorithms resemble ASIC algorithms, but the constraints are those of a finished chip.

## LUTs as truth tables

The basic FPGA logic element is a $k$-input LUT: a small memory with $2^k$ configuration bits, one per row of the truth table, so it can implement any function of $k$ inputs. A 6-input LUT holds 64 bits. Technology mapping for FPGAs therefore covers the network with $k$-feasible cones (subgraphs with at most $k$ inputs) rather than with library cells, and the cost is the number of LUTs and the depth of the LUT network. Algorithms such as FlowMap find depth-optimal mappings in polynomial time, which has no direct ASIC analogue. Hard blocks such as block RAMs, DSP multipliers and transceivers are inferred where the RTL allows.

## Placement and routing on a fixed fabric

FPGA placement assigns packed logic blocks to legal sites of matching type on a grid, and the objective combines wirelength and timing much as in ASIC placement, with simulated annealing (as in VPR) or analytical methods. Routing is different: the wires already exist, so the router selects a path through a routing-resource graph of wire segments and programmable switches. Two nets cannot share a wire, so routers such as PathFinder negotiate congestion by raising the cost of overused resources over iterations until every net has its own path. Timing comes from the fixed delays of the chosen segments.

## Vendor tools, open tools and trade-offs

Commercial FPGA flows come from the device vendors: AMD Vivado for Xilinx devices and Quartus for Altera devices. Open flows combine Yosys for synthesis with nextpnr or VPR for placement and routing, using documented bitstream formats for some device families. Because the fabric is fixed, an FPGA implementation runs out of resources (LUTs, block RAM, routing) rather than area, and timing closure cannot add wires, only choose them. The flow runs in minutes to hours, so FPGA designers iterate far more often than ASIC designers. Timing closure on an FPGA therefore means choosing better placements and routes, adding pipeline registers, or moving logic into hard blocks.

## Going deeper

Betz, Rose and Marquardt, "Architecture and CAD for Deep-Submicron FPGAs" (1999); Cong and Ding, "FlowMap", IEEE Transactions on CAD (1994); McMurchie and Ebeling, "PathFinder" (FPGA 1995).
