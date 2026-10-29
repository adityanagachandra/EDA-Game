# SDC timing constraints

Synopsys Design Constraints (SDC) is the Tcl-based format in which designers state timing intent: which clocks exist and at what period, when inputs arrive and outputs are required relative to those clocks, and which paths are exceptions. Synthesis, placement, CTS, routing optimization and static timing analysis all read the same SDC. A wrong constraint is worse than a slow tool: every engine will optimize faithfully toward the wrong target.

## Clocks

The command `create_clock -name clk -period 2.5 [get_ports clk]` declares a clock with a 2.5 ns period, which is $1000/2.5 = 400$ MHz. Generated clocks derived from it, for example a divide-by-two, are declared with `create_generated_clock`. `set_clock_uncertainty` subtracts margin for jitter and expected skew from every check, and `set_clock_latency` describes clock delay before CTS has built the real tree. Clocks with no defined relationship are declared asynchronous so that the timer does not check impossible paths between them. The clock definitions are the root of every timing check, so a wrong period or a missing clock makes every downstream slack meaningless.

## Input and output delays

A block's ports connect to logic outside it, so SDC states how much of the clock period that outside logic uses. `set_input_delay 0.8 -clock clk [get_ports din]` says data reaches the port 0.8 ns after the clock edge. For a register inside the block, the logic from that port must settle by the next edge minus setup and uncertainty. With a 2.5 ns period, a 0.8 ns input delay, a 0.1 ns setup time and 0.15 ns of uncertainty, the internal path has $2.5 - 0.8 - 0.1 - 0.15 = 1.45$ ns. Output delays are the mirror: `set_output_delay` reserves time for logic after the output port.

## Exceptions

Some paths should not be checked with the default single-cycle rule. `set_false_path` removes paths that can never be sensitized or that cross between asynchronous clocks. `set_multicycle_path 2 -setup` lets a path take two cycles, so its setup check moves to the second capture edge: with a 2.5 ns clock, the path has 5 ns minus setup and uncertainty. Exceptions are powerful and dangerous. A false-path constraint written too broadly hides real failures, and nothing in the timer will report it, which is why constraint review and constraint linting are signoff steps.

## Going deeper

Gangadharan and Churiwala, "Constraining Designs for Synthesis and Timing Analysis: A Practical Guide to Synopsys Design Constraints (SDC)" (2013).
