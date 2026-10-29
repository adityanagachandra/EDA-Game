# SystemVerilog always_ff and always_comb

SystemVerilog adds always_ff, always_comb and always_latch: procedural blocks that state the designer's intent so that tools can check it. They simulate almost like the Verilog always blocks they replace, but they close several simulation-synthesis gaps and let a linter or synthesizer reject code that would not produce the intended hardware.

## Intent that tools check

always_ff declares a clocked block that must infer flip-flops; always_comb declares combinational logic; always_latch declares an intentional latch. A tool that finds a latch inferred inside always_comb must warn or error, because the code contradicts the stated intent. Without the keyword, the same latch would pass silently as a plain always block.

Each of the three also claims exclusive ownership of the variables it writes: a variable written in always_comb may not be written by any other process, including a continuous assignment. This rules out multiple-driver bugs at compile time. The keywords therefore turn coding guidelines into checked rules, which is exactly what a compiler should do with intent.

## always_comb semantics

always_comb is more than always @(*). Its implicit sensitivity list includes variables read inside functions it calls, where @(*) looks only at the block's own expressions; a change to a module-level variable read in a called function wakes always_comb but not @(*). It also runs once automatically at time zero, after initial and always blocks start, so its outputs are consistent from the beginning even if no input ever changes.

Variables written in the block are excluded from its sensitivity, so the block does not retrigger itself. A temporary written before it is read becomes a wire carrying an intermediate value, exactly as with blocking assignments in a plain combinational block. These rules make simulation match what synthesis builds.

## Clocks, types and what remains unchecked

always_ff @(posedge clk or negedge rst_n) lists two edges, and the tool must identify one as the clock and the other as an asynchronous control, from the if (!rst_n) ... else ... pattern. SystemVerilog code typically uses the four-state type logic for everything these blocks write, which removes Verilog's reg-versus-wire choice; nets with several drivers still need wire.

Some problems remain beyond the keywords. A blocking assignment inside always_ff to a variable read by another block is legal but race-prone, so linters flag it. A combinational loop through two always_comb blocks is not detected by either block's checks; it needs a loop check across the netlist. Knowing the limits of a check is part of building on it.

## Going deeper

IEEE Std 1800-2017, clause 9.2.2; Sutherland, "SystemVerilog for Design" (2nd ed.), chapter 6.
