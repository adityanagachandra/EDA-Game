# SDC timing constraints

A timer knows nothing about intent until it reads constraints. The Synopsys Design Constraints format, SDC, is a set of Tcl commands that define clocks, the environment at the block boundary, and the exceptions to single-cycle timing. Every slack in a report is only as correct as the SDC behind it, and constraint bugs are among the most common causes of silicon timing failures.

## Defining clocks

create_clock names a clock, its source port or pin, its period and its waveform. create_generated_clock derives a clock from a master, by -divide_by, -multiply_by or -edges, and keeps the phase relationship so paths between the two domains are timed correctly. set_clock_latency gives latency: -source for the delay from the clock's origin to its definition point, and without -source an ideal network latency used before clock tree synthesis. set_propagated_clock tells the timer to compute network latency from the built tree instead. set_clock_uncertainty reserves margin for jitter and, before CTS, expected skew, with separate -setup and -hold values. Declaring a divided clock as an independent create_clock is a typical mistake: it severs the relationship to the master.

## Constraining the block boundary

Inside a block the timer sees flops; at its ports it needs a model of the outside. set_input_delay gives the arrival at an input port relative to a clock edge, standing in for an external launch flop and its logic. set_output_delay gives how much of the cycle the logic and setup outside the block consume after an output port, so the port's required time is the capture edge minus that value. set_driving_cell or set_input_transition gives input ports a realistic slew, and set_load adds the capacitance an output port drives. An input port with no input delay is unconstrained: its paths are never reported, which is why constraint checks flag missing I/O constraints.

## Object queries and units

SDC commands apply to design objects found with get_ports, get_pins, get_cells, get_nets and get_clocks. The query matters because one name can denote different objects, and a query that matches nothing makes the constraint silently disappear. SDC numbers carry no units: the timer interprets them in the time unit of the library or session, usually ns, so a team that writes 200 meaning ps against an ns library sets a 200 ns constraint. Robust flows check that every clock reaches its sinks, every port is constrained and every exception matches something, using reports such as check_timing. In OpenSTA, read_sdc applies the file after the design is linked, so the objects it names must exist.

## Going deeper

The Synopsys Design Constraints format specification; Gangadharan and Churiwala, "Constraining Designs for Synthesis and Timing Analysis: A Practical Guide to Synopsys Design Constraints (SDC)" (2013).
