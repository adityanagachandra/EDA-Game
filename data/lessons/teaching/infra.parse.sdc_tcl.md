# SDC as Tcl

Synopsys Design Constraints (SDC) is the language of timing intent: clocks, input and output delays, exceptions. It is not a standalone grammar but a set of Tcl commands, so an SDC file is a Tcl program. A tool that reads SDC must evaluate Tcl and must resolve object queries against its own database. Knowing this explains how constraints work, and how they silently break.

## SDC is a Tcl dialect

Each SDC line is a Tcl command: `create_clock -name clk -period 2.5 [get_ports clk]`. Tcl rules apply. Square brackets are command substitution: `get_ports clk` runs first and its result becomes an argument. A dollar sign before a name substitutes the variable's value, so `set p 2.5` followed by a reference to p after -period is legal. Braces quote literally, so `{in_a in_b}` is one word holding a two-element list. Loops, procedures and expressions are all available: a foreach loop over {a b c} whose body calls set_load on get_ports of the loop variable issues three constraints. A tool that only pattern-matches SDC lines will misread any file that uses these features, which is why real readers embed a Tcl interpreter and define the SDC commands inside it.

## Object queries return collections

Commands such as `get_ports`, `get_pins`, `get_cells`, `get_nets` and `get_clocks` look up database objects by name pattern and return a collection of object handles, not strings. `get_pins u1/*` returns every pin of u1. The constraint command then attaches to those objects. Two consequences follow. A pattern that matches nothing returns an empty collection, and the constraint silently applies to nothing, so careful tools warn on empty matches. And constraints bind through names: if synthesis renames `u_alu/add_0` to `u_alu/add_0_inst`, a constraint written against the old hierarchical name stops matching. Collections exist because re-resolving names for every command would be slow and ambiguous when a port and a net share a name.

## Interpreting SDC inside a tool

SDC numbers are unitless. `set_input_delay 300 -clock clk [get_ports a]` means 300 of whatever time unit is in effect, set by `set_units` or inherited from the library, so a file written for picoseconds and read in nanoseconds is off by a thousand. Tools differ in how much Tcl they allow. OpenSTA defines the SDC commands in its own Tcl interpreter and sources the file. Some sign-off flows use restricted readers that reject arbitrary Tcl, so that a constraint file cannot run side effects and means the same in every tool. Writing SDC back out after optimization is the reverse problem: collections of handles must be turned into names that resolve in the modified netlist.

## Going deeper

The SDC specification published with the open-source SDC parser; the OpenSTA documentation for read_sdc; Gangadharan and Churiwala, "Constraining Designs for Synthesis and Timing Analysis" (2013).
