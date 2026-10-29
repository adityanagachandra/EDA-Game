# The SKY130 open PDK

SKY130 is the SkyWater 130 nm process, released as an open process design kit with Google's backing. It was the first open PDK that a full digital flow could target end to end, so it is where most open-source EDA tools are tested and compared. For a tool builder it is a real process to read, parse and run against: every view a commercial flow consumes exists here in open form.

## What the SKY130 kit contains

A PDK is the process described in files. SKY130 ships SPICE device models (for ngspice and Xyce), the layer definitions and design rules, DRC and LVS decks for open verifiers such as Magic, KLayout and Netgen, and libraries built on the process: standard cells, IO cells and some primitive devices. Each library cell comes in several views: a GDS layout, a LEF abstract, Liberty timing and power models per corner, a SPICE netlist and a Verilog model. What the PDK does not contain is anything specific to your design: RTL, SDC, floorplan choices and flow parameters are yours. When two people disagree about a DRC or timing result on the same layout, the first thing to compare is the PDK version, because rule decks, models and libraries evolve between releases. A flow that does not pin its PDK version is not reproducible.

## The SKY130 layer stack and libraries

SKY130 has a local-interconnect layer, li1, below five metal layers named met1 to met5. Standard cells use li1 for their internal wiring, which leaves more of met1 free for routing. Higher layers are thicker and less resistive, so they suit long wires and the power grid. Several standard-cell libraries share the process, named like sky130_fd_sc_hd: fd marking a foundry-provided library, sc for standard cell, and a suffix for the variant. hd is high density, hs is high speed, and there are others for low power, low leakage and different voltages. They trade cell height, transistor width and leakage, so the same RTL maps to different area and speed. A worked choice: a small controller with a relaxed 50 MHz clock and a tight area budget goes to hd. A datapath that misses timing at hd may close in hs, at the price of more area and leakage.

## How a flow consumes SKY130

Each flow stage reads the view it needs. Synthesis (Yosys and ABC) reads Liberty to know cell functions, area and delay. Placement and routing read the technology LEF (layers, vias, sites) and the cell LEF (size, pins, obstructions); they never touch transistor geometry. The static timer reads Liberty for the chosen corner, and extraction uses the layer stack to compute parasitics. At the end the flow replaces each LEF abstract with the cell's GDS (ORFS does this with KLayout), and the DRC and LVS decks check the merged layout. Two common mistakes come straight from this division. A final GDS with empty boxes where cells should be means the merge step did not run. Timing that looks too good often means a fast-corner Liberty was loaded, for example an ff file instead of tt, so every arc is faster than the corner you meant to sign off.

## Going deeper

The SkyWater SKY130 PDK documentation; the open_pdks installer documentation; the OpenROAD-flow-scripts platform files for sky130hd.
