# Blocking and nonblocking assignments

Verilog has two procedural assignments. Blocking assignment (=) updates its variable immediately; nonblocking assignment (<=) samples its right-hand side now and updates after the block has run. The difference decides whether a clocked block describes a shift register or a single register, and it is the most common source of RTL bugs. Every simulator and synthesizer must implement the two exactly, so this concept is computed, not memorized.

## Nonblocking: sample now, update later

A nonblocking assignment evaluates its right-hand side in the active region and schedules the update for the NBA region. Within one clock edge, every statement therefore reads the values from before the edge. Take a = 1, b = 2, c = 3 and the block

always @(posedge clk) begin a <= b; b <= c; c <= a; end

After one edge, a = 2, b = 3 and c = 1: the three values rotate, regardless of statement order. After a second edge they rotate again to a = 3, b = 1, c = 2. This is exactly what three flip-flops wired in a ring do, which is why nonblocking assignment is the rule for clocked logic.

## Blocking: update immediately

A blocking assignment updates its variable before the next statement runs, so later statements see the new value. With the same starting values, the block a = b; b = c; c = a; gives a = 2, b = 3 and then c = a = 2: the first value is lost. Order now matters. Reversing the statements to c = a; b = c; a = b; gives c = 1, b = 1, a = 1.

Inside combinational blocks this immediacy is what you want: blocking assignments describe a sequence of intermediate values, exactly like variables in software, and synthesis builds the corresponding chain of logic. Compilers make this precise with static single assignment: each blocking write creates a new version of the variable, and each read refers to the most recent version above it.

## Mixing the two in one block

In a mixed block, execute statements in order, apply blocking updates immediately and collect nonblocking updates until the end. With a = 2, b = 9, c = 8 and the block b <= c + 2; a = c + a; c = a + 4;, the first statement schedules b = 10; the second sets a = 10; the third reads the new a and sets c = 14. After the edge b becomes 10.

Mixing is legal but dangerous. A blocking write to a variable that another clocked block reads reintroduces races, and lint tools warn about blocking and nonblocking writes to the same variable. A common safe idiom does use both: blocking writes to a local temporary that is written before it is read, followed by nonblocking writes of the registers, which models intermediate logic feeding flip-flops.

## Going deeper

Cummings, "Nonblocking Assignments in Verilog Synthesis, Coding Styles That Kill!" (SNUG 2000); IEEE Std 1800-2017, clauses 10.4.1 and 10.4.2.
