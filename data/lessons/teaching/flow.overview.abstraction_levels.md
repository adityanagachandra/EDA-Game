# Abstraction levels from system to layout

Designers describe a chip at several levels of abstraction: system, algorithm or behavior, register transfer, gate, transistor and layout. Each level hides detail that the level below fills in. EDA tools are mostly translators and checkers between adjacent levels, so knowing what each level keeps and drops tells you what a tool can compute from it and what it must assume.

## The ladder of levels

At the system level a design is a set of communicating components (processors, memories, accelerators, interconnect) with throughput and latency budgets. At the behavioral level it is an algorithm, often in C++ or SystemC, with no clock. At the register-transfer level (RTL) it is registers and the combinational logic between them, cycle by cycle, written in Verilog, SystemVerilog or VHDL. At the gate level it is a netlist of library cells such as NAND2 and DFF. At the transistor level it is a circuit of MOSFETs, resistors and capacitors. At the layout level it is polygons on mask layers. Gajski and Kuhn's Y-chart organizes the same idea along three views (behavioral, structural, physical) with these levels as rings.

## What each level keeps and drops

Going down the ladder adds information. RTL fixes the cycle on which each value is computed but not which gates compute it. The gate netlist fixes the gates but not where they sit. Layout fixes geometry, and with it the wire resistance and capacitance that set real delay. Going up hides information, which is what makes large designs tractable: RTL simulation runs much faster than transistor simulation because it ignores voltages and computes one value per signal per cycle. The price is that a property that depends on dropped detail (exact delay, crosstalk, IR drop) cannot be computed at the higher level, only estimated.

## Refinement and checking between levels

Every step down is a refinement, and every refinement can introduce an error. So the flow pairs each transformation with a check that the levels agree. Logic synthesis is checked by logic equivalence checking between RTL and netlist. Layout is checked by layout-versus-schematic (LVS), which extracts a transistor netlist from polygons and compares it with the intended one. Timing estimated at a high level is re-checked once real parasitics are known. When a check fails, the designer goes back up a level. This is the first reason flows iterate. Each check is itself an EDA engine, so the checks scale with the design just as the transformations do.

## Going deeper

Gajski and Kuhn, "Guest Editors' Introduction: New VLSI Tools", IEEE Computer (1983), which introduced the Y-chart.
