# The RTL-to-GDSII flow end to end

The RTL-to-GDSII flow is the chain of tools that turns register-transfer code into the mask layout a foundry manufactures. It is the backbone of digital chip design and of this curriculum: almost every later domain is one stage of it, studied in depth. Knowing the order of the stages, and what each one hands the next, is how a tool designer knows what their engine can assume about its input.

## The stages in order

A digital flow runs: specification and architecture, RTL design, functional verification, logic synthesis, design-for-test (scan) insertion, floorplanning, placement, clock tree synthesis, routing, parasitic extraction, signoff (static timing, power and IR drop, design-rule and layout-versus-schematic checks), tape-out, then fabrication, packaging and production test. Verification is not a single box: it runs alongside RTL design and continues on the netlist. The front end is everything up to a verified netlist; the back end, or physical design, is floorplanning through signoff. Real flows add steps such as power-intent checks, timing-driven optimization after each physical stage, and engineering change orders late in signoff.

## What each stage hands the next

Each stage consumes one representation and produces a more detailed one. Synthesis reads RTL, constraints (SDC) and the Liberty library, and writes a gate netlist. Floorplanning adds a die outline, macro positions and the power grid. Placement adds a legal location for every standard cell. Clock tree synthesis adds buffers and wires that deliver the clock to every flip-flop with controlled skew. Routing adds the geometry of every wire and via. Extraction turns that geometry into resistances and capacitances (SPEF). Signoff reads all of it and either passes the design or reports violations. The final layout is streamed out as GDSII or OASIS.

## What becomes known at each stage

The flow is also a schedule of when information appears. Gate count and logic depth are first known after synthesis. Cell positions first exist after placement, so wirelength before placement is only an estimate. Clock skew is first real after clock tree synthesis. Wire geometry and coupling capacitance first exist after routing and extraction. Every earlier stage must therefore optimize against estimates of quantities that a later stage will compute exactly. When the estimate was wrong, a later stage fails and the flow loops back, which is the theme of design closure.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 1; the OpenROAD flow scripts, which implement the whole chain with open-source tools.
