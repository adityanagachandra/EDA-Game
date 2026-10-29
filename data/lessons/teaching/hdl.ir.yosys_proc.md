# The Yosys proc passes

The proc passes are where Yosys turns procedural HDL code into a netlist. Every always block arrives as an RTLIL process, and proc converts its decision tree into multiplexers and its clock and reset structure into flip-flops and latches. Latch inference, reset recognition and the shape of the initial netlist are all decided here, which makes proc the first place to look when synthesis output surprises you.

## From processes to cells

A process holds a decision tree, built from the if and case statements, and sync rules that say when assignments take effect (on a clock edge, or always for combinational code). proc lowers it to cells: each signal's final value through the decision tree becomes a tree of multiplexer cells, and each sync rule on an edge becomes a flip-flop whose D input is that multiplexer tree. An always @(posedge clk) process therefore becomes a dff cell fed by muxes.

Running proc early matters, because the optimization passes (opt_expr, opt_merge, opt_clean) work on cells and nets, and a process is still procedural. Every standard synthesis script runs proc right after reading and checking the hierarchy.

## The sub-passes

proc is a sequence of smaller passes. proc_clean and proc_rmdead remove empty parts of processes and decision-tree branches that can never be taken. proc_init turns initial values into init attributes on wires, which FPGA flows honor. proc_arst looks for an asynchronous reset: a top-level if on a signal that also appears as an edge in the sensitivity list, and rewrites it so that a flip-flop with asynchronous reset can be built. proc_mux converts the decision trees into multiplexer cells. proc_dlatch creates latches for signals a combinational process does not assign on every path. proc_dff creates flip-flops for the edge-triggered sync rules.

Each step has a narrow job, and the order matters: reset detection must happen before the flip-flops are built, and latch creation needs the multiplexer structure to know which paths leave a signal unassigned.

## Reading what proc produced

A latch cell in the log after proc means some output of a block meant to be combinational is not assigned on every path; proc_dlatch only creates one when a value must hold. A missing reset pin on a flip-flop that should have one usually means proc_arst did not recognize an unusual reset coding style, so the reset ended up in the data path.

Blocking assignments to temporaries are tracked through the decision tree, so later reads use the newly assigned value, much like SSA renaming; the temporary becomes plain logic unless it is read before it is written, in which case it is a flip-flop. The multiplexer trees proc builds are a direct translation and often redundant, with branches assigning the same value or controlled by constants, so the opt passes that follow (opt_muxtree, opt_expr, opt_clean) usually shrink them a lot.

## Going deeper

The Yosys manual, the proc command and its sub-passes; Wolf, "Yosys Open SYnthesis Suite" documentation, chapter on the synthesis flow.
