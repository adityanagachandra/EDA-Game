# The open analog flow: xschem, ngspice, Magic, Netgen

Open PDKs such as SkyWater SKY130 made a complete open-source analog flow practical: schematic capture, simulation, layout, DRC, extraction and LVS, all with inspectable tools and file formats. For an analog EDA tool designer this flow is a testbed: a new sizer, layout generator or agent can be inserted at any stage and checked against real signoff steps.

## The tools and what each does

xschem is the schematic editor: it captures the circuit with PDK symbols and writes a SPICE netlist. ngspice simulates that netlist with the PDK's device models, running DC, AC, transient, noise and Monte Carlo analyses; Xyce is an alternative engine for large or parallel runs. Magic is the layout editor: it draws polygons on the PDK's layers, checks design rules interactively and extracts a netlist with devices and parasitic capacitances from the layout. Netgen performs layout-versus-schematic comparison between the extracted netlist and the schematic netlist. KLayout is often used alongside for viewing, scripting and additional DRC decks. Together they cover the custom analog loop that commercial suites cover with their own tools.

## Signoff steps in order

The order matters. First the layout must be DRC clean, because extraction of an illegal layout is meaningless. Then LVS must pass: Netgen treats both netlists as graphs and searches for an isomorphism by iterative partition refinement, since extracted node names are arbitrary. It reduces parallel and series devices first, so two parallel 1 micrometre devices match one 2 micrometre device if the rules allow it. Only then is parasitic extraction trusted, and the extracted netlist, including coupling capacitances, is simulated against the specs. A layout that passes DRC and LVS but fails post-layout simulation is still a failure. Each step produces a report, and those reports, not tool exit codes, are the evidence of correctness.

## The flow as a testbed for new tools

Because every engine and format is open, the flow can be driven by scripts: xschem and Magic run in batch mode, ngspice reads control scripts, Netgen reads setup files. A new automatic layout tool can emit GDS or Magic layout and be judged by the same DRC, LVS and post-layout simulation a manual layout faces. A sizing optimizer can rewrite netlist parameters and call ngspice in a loop. An agent can orchestrate the whole chain. The PDK enters at every stage, as device models, symbols, design rules, extraction rules and LVS setup, so results transfer only within that process. For research and teaching the flow's value is that every number can be traced to an inspectable tool, which is what a tool designer needs to debug their own algorithm.

## Going deeper

The documentation of xschem, ngspice, Magic and Netgen; the SkyWater SKY130 open PDK documentation.
