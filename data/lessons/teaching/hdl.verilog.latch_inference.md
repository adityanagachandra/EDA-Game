# Inferred latches

A combinational always block that does not assign a variable on every path forces synthesis to remember the old value on the missing paths, and the only hardware that remembers without a clock is a latch. Inferred latches are among the most common RTL defects: they complicate timing analysis, break scan insertion and often reflect a forgotten else. Counting them precisely is what a lint or synthesis tool does.

## Incomplete if statements

In always @(*), a variable must receive a value on every path through the block. if (s) y = a; with no else leaves y unassigned when s is 0, so its value must hold: synthesis builds a level-sensitive latch enabled by s. Each such variable contributes its full width in latch bits.

The cheapest cure is a default assignment at the top of the block, such as y = 0;, placed before the if. Then every path assigns y, the last assignment wins, and synthesis builds a multiplexer. For example, with y1 (8 bits) and y2 (4 bits) assigned only inside if (s) and y3 (2 bits) defaulted at the top, the block infers $8 + 4 = 12$ latch bits.

## If/else with uneven assignments

An if with an else covers both values of the condition, but only if each variable is assigned in both arms. With if (s) begin y1 = a; y3 = b; end else begin y3 = a; end, y3 is complete and becomes a multiplexer, while y1 is assigned only when s is 1 and becomes a latch. The rule per variable is: no latch if the variable has a default, or if it is assigned in every arm.

This is a reaching-definitions question in disguise. At the end of the block, ask whether every path delivers a definition of the variable. If some path delivers only the value from before the block, that value must be stored. Compilers use the same analysis, and the Yosys pass that creates the latch, proc_dlatch, runs exactly at this point.

## Case statements and full coverage

A case statement is complete for a variable only if the variable is assigned in every item and the items cover every selector value. With a 2-bit selector, four items 2'd0 to 2'd3 cover all values, so a variable assigned in all four needs no latch. With only three items, the fourth value leaves every variable unassigned, so each one latches unless a default item assigns it or a default assignment precedes the case.

Two traps deserve names. First, synthesis pragmas such as full_case tell synthesis to treat missing values as don't-cares and avoid the latch, but simulation still holds the old value, which creates a simulation-synthesis mismatch. Second, a default item assigns only the variables written in it; a variable missing from the default still latches. SystemVerilog's always_comb asks tools to report any such latch.

## Going deeper

The Yosys manual, the proc_dlatch pass; Mills and Cummings, "RTL Coding Styles That Yield Simulation and Synthesis Mismatches" (SNUG 1999).
