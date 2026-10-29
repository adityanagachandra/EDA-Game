# Memory macros and non-volatile memory

Memories are not synthesized gate by gate. They arrive as macros produced by a memory compiler: a layout, a timing model, an abstract placement view and a behavioral model, all for a specific configuration. Non-volatile memories (flash, and emerging resistive and magnetic types) extend the family when data must survive power-down. A tool designer needs to know what each view promises and what each technology trades.

## What a memory compiler delivers

Given words, bits, column-mux factor, banking and options such as write masks or redundancy, a memory compiler emits a hard macro with a set of views: the layout (GDSII or OASIS), an abstract physical view (LEF) with pin locations and blockages for placement and routing, a timing and power model (Liberty) characterized at each corner, and a behavioral Verilog model for simulation.

Implementation tools never look inside the macro; they place it as a block, route to its pins and time it through its Liberty arcs (clock-to-output for reads, setup and hold on address and data inputs). The macro's correctness inside is the compiler vendor's responsibility, verified by its own characterization.

## Test and repair

Memories have the densest layout on the chip and therefore the most defects. They are tested by memory built-in self-test (MBIST): an on-chip engine runs march algorithms (sequences of writes and reads in ascending and descending address order) that detect stuck-at, transition and coupling faults.

Large memories also carry redundant rows and columns. When MBIST finds a failing row, built-in self-repair or fuse programming remaps it to a spare, raising yield substantially. The redundancy, the repair logic and the BIST controller are all inserted by the flow, which is why memory macros come with their own test views and why design-for-test tools must understand them.

## Non-volatile memories

Embedded flash stores charge on a floating gate or in a charge-trap layer; it retains data for years without power, but writes are slow, need high voltage and wear the cell out after a limited number of program and erase cycles, and it needs extra process steps. Emerging non-volatile memories store state differently: MRAM in the magnetization of a magnetic tunnel junction, ReRAM in a resistive filament, phase-change memory in the crystalline or amorphous state of a chalcogenide. They write faster than flash and integrate in the back-end metal layers, at the price of variability and sensing difficulty.

Read-only memory (mask ROM) is fixed at fabrication. For a designer the choice is a matrix of density, endurance, write energy, retention, process cost and read speed, and the macro's views (timing, power, test) differ accordingly.

## Going deeper

van de Goor, "Testing Semiconductor Memories: Theory and Practice" (1991); Yu and Chen, "Emerging Memory Technologies" (IEEE Solid-State Circuits Magazine, 2016).
