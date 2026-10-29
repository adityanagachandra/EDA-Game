# The synthesizable subset

Verilog was designed for simulation, and only part of it describes hardware. The synthesizable subset is the part that logic synthesis maps to gates, registers and memories; the rest serves testbenches. A tool builder must know where the boundary lies, what synthesis silently ignores, and which constructs synthesize but behave differently in simulation, because those mismatches are where real bugs hide.

## Constructs synthesis ignores or rejects

Delays such as assign #5 y = a & b; are ignored by synthesis, with a warning: hardware timing comes from cells and wires, not from annotations. Printing system tasks and file I/O are skipped. General event waits inside procedural code, fork and join, and forever loops that toggle clocks describe testbench sequencing and are rejected.

Initial blocks depend on the target. FPGA configuration can load initial register values, so FPGA synthesis honors simple initial assignments; ASIC flip-flops power up unknown, so ASIC flows ignore them and state must come from reset logic. Writing RTL that relies on initial values therefore works on an FPGA and fails on an ASIC, a portability trap worth knowing when building flows for both.

## Static loops, functions and memories

A loop is synthesizable when its trip count is known at compile time. for (i = 0; i < 8; i = i + 1) y[i] = a[7 - i]; unrolls into eight wire connections, and the 32-bit integer loop variable disappears entirely, because each unrolled copy sees a constant. A loop bounded by an input port describes no fixed hardware and is rejected.

Functions without timing controls are inlined as combinational logic at each call. Arrays such as reg [7:0] mem [0:255] are recognized by memory inference: if the read and write pattern matches a RAM the target provides, synthesis maps the array to a block RAM or macro; otherwise it falls back to flip-flops and multiplexers, which for large arrays is enormous. Writing memories in the pattern a tool recognizes is a real design skill.

## Simulation-synthesis mismatches

The dangerous constructs are those that synthesize but mean something else in simulation. An incomplete sensitivity list makes simulation miss changes that the synthesized gate responds to. A full_case pragma tells synthesis that unlisted case values are don't-cares while simulation holds the outputs. Assigning y = 1'bx in an unused branch lets synthesis choose any value that minimizes logic while simulation shows x. The === operator compares x literally, which no hardware can do.

The defensive style follows directly: use always @(*) or always_comb, give every combinational output a default at the top of the block, prefer SystemVerilog unique and priority over pragmas, and verify the netlist against the RTL with an equivalence checker. Each rule removes one gap between the two models a tool chain maintains.

## Going deeper

IEEE Std 1364.1-2002 (Verilog register-transfer level synthesis); Mills and Cummings, "RTL Coding Styles That Yield Simulation and Synthesis Mismatches" (SNUG 1999).
