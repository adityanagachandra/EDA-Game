# OpenLane 2

OpenLane 2 is a Python infrastructure for building RTL-to-GDS flows out of steps, the successor of the original OpenLane. It does not contain its own placer or router; it orchestrates Yosys, OpenROAD, KLayout, Magic, Netgen and other tools on open PDKs. For a tool builder it is a second, differently shaped harness next to ORFS, and a model for how to make a flow extensible.

## Flows built from steps

In OpenLane 2 a step wraps one action, such as running synthesis, global placement or a DRC check. Each step declares the design views it consumes and produces: netlists, ODB databases, DEF, SPEF, GDS and reports. A flow is a sequence of steps, and the state of views passes from one to the next, so a later step starts where the earlier ones left off. Flows are written in Python, which makes it natural to add a custom step, for example a checker that inspects placement density before clock tree synthesis, or to build a new flow that reuses existing steps in a different order. The configuration supplies the design's files and parameters, and per-PDK defaults fill in the rest. This structure is the flow-level version of a plug-in architecture: well-defined interfaces between steps make the pieces replaceable.

## How OpenLane 2 relates to OpenROAD, Yosys and PDKs

OpenLane 2 calls the same engines ORFS calls. Synthesis steps run Yosys with ABC. Floorplanning, placement, clock tree synthesis, resizing, routing, extraction and timing run OpenROAD's tools on OpenDB. Physical verification runs Magic or KLayout with the PDK's DRC decks, and Netgen compares the extracted layout netlist with the gate netlist for LVS. The PDK, for example SKY130 or GF180MCU, supplies the views and the per-PDK settings; OpenLane holds no process technology of its own. A consequence for debugging: when a routing step produces bad wires, the code responsible is OpenROAD's router, and OpenLane's role is how it invoked it, which parameters and inputs it passed. Reproducing the step's tool invocation outside the flow is the fastest way to tell an orchestration bug from an engine bug.

## Choosing between OpenLane 2 and ORFS

Both are orchestration layers over the same engines, so the choice is about workflow. ORFS is developed alongside OpenROAD, is driven by Makefiles and Tcl stage scripts, and exposes OpenROAD's commands directly, which suits engine development and reproducing OpenROAD results. OpenLane 2 offers a Python API for composing flows from steps, plus configuration conventions and checks, which suits building custom flows and integrating new checks. On the same design and PDK, the two may report different QoR: default utilization, density, repair settings, routing layers and signoff checks differ, so the recipes differ even though the engines do not. A fair comparison fixes those knobs to the same values, or reports each flow's settings beside its numbers.

## Going deeper

The OpenLane 2 documentation (steps, flows and configuration); the OpenROAD-flow-scripts documentation.
