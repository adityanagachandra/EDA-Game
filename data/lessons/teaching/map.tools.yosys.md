# Yosys techmap, dfflibmap and abc

Yosys turns RTL into a netlist and hands combinational mapping to ABC. In an ASIC flow the mapping steps are a short sequence of passes, each with a distinct job: lower generic cells, map flip-flops to the library, map combinational logic with ABC, then clean up. Knowing which pass owns which decision tells you where to look when a mapped netlist is wrong or poor.

## The standard-cell mapping sequence

After synth has produced Yosys's internal gate-level cells, a typical ASIC script runs dfflibmap -liberty cells.lib to replace internal flip-flops with library flip-flops, then abc -liberty cells.lib to map the combinational logic, then opt_clean, and finally stat -liberty cells.lib to report area. read_liberty -lib cells.lib, run early, imports the cells as black boxes so later passes recognize instances of them. dfflibmap runs before abc so that ABC sees flip-flops as fixed cells and maps only the logic between them. ABC itself never maps sequential elements in this flow.

## techmap and the internal cell library

techmap replaces cells by implementations given in a Verilog map file; with no file it lowers high-level cells, such as adders and comparators, into Yosys's internal gate cells. It is rule-based substitution, not optimization: each cell is rewritten independently. Custom map files are how a flow targets special structures, such as FPGA primitives, carry chains or hand-designed macros. hilomap maps constant drivers onto tie-high and tie-low cells. Passes such as dfflegalize convert flip-flop types the library lacks, for example one with asynchronous reset, into supported types with extra logic.

## Options that steer ABC

abc -lut k maps into k-input LUTs instead of cells; abc -D sets a delay target in picoseconds, which Yosys substitutes into its ABC script for timing-driven mapping; abc -script supplies a custom ABC script, which is where structural choices or extra area recovery are added. Yosys flows for FPGAs often use abc9, which passes timing and box information for hard blocks to ABC. Whatever the options, the mapped result should be checked with an independent equivalence check against the pre-mapping design, since Yosys and ABC each transform the netlist.

## Going deeper

Wolf, "Yosys Open SYnthesis Suite" manual and the help text of the techmap, dfflibmap, abc and stat passes; the OpenROAD flow scripts, which use this sequence for ASIC synthesis.
