# Flip-flop inference from clocked blocks

Synthesis does not see flip-flops in RTL; it infers them. The rule is about lifetimes: a variable needs storage when its value must survive from one clock cycle to the next. Counting inferred flip-flop bits is a direct test of whether you read a clocked block the way a synthesizer does, and getting it wrong is how designs gain phantom registers or lose intended ones.

## Every nonblocking target is a register

In always @(posedge clk), each variable assigned with <= becomes a flip-flop as wide as the variable. A conditional assignment such as if (en) q <= d; is still a register; on cycles where en is low it keeps its value, so synthesis builds a flip-flop with an enable (or a feedback multiplexer). A reset branch adds a reset or set to the same flip-flop; it does not add a second one. Continuous assignments, by contrast, never infer storage.

Count bits, not variables. With reg [7:0] q1 and reg [11:0] q2 both assigned with <= in a clocked block, and a 4-bit continuous assign elsewhere, synthesis infers $8 + 12 = 20$ flip-flop bits before optimization. Later passes may remove some (a register with a constant input, or one whose output is unused), which is why the count is stated before optimization.

## Blocking temporaries written before read

A clocked block may also contain blocking assignments to temporaries: t = a + b; q <= t;. Here t is written before it is read on every execution of the block, so no read ever sees a value from a previous cycle. Its value lives only within the edge, and synthesis turns t into a wire carrying the sum. No flip-flop is built for t, provided nothing outside the block reads it.

The test is the same one a compiler applies with liveness: is the variable live on entry to the block? A written-before-read temporary is dead on entry, because every path overwrites it first. Adding the temporary changes the flop count by zero. In the example q1 (1 bit), q2 (4), q3 (16) and an 8-bit temporary written before read, the count is $1 + 4 + 16 = 21$.

## Temporaries read before written

Reverse the order and the meaning changes: q1 <= t1; t1 = b - a;. Now the read of t1 sees the value written on the previous clock edge, so t1 must be stored across cycles. Synthesis infers a flip-flop for t1, and q1 receives a twice-delayed difference: the hardware is a two-stage pipeline although only one nonblocking assignment appears.

This is a classic accidental register. The rule generalizes: a variable assigned in a clocked block needs storage if some read of it, inside or outside the block, can observe a value from an earlier cycle. Reads before writes, reads on paths where the write is skipped, and reads from other blocks all count. Lint tools warn about blocking assignments in clocked blocks for this reason, and a counting exercise such as q1 (3 bits), q2 (4), a 16-bit temporary read before written and a 16-bit temporary written before read gives $3 + 4 + 16 = 23$.

## Going deeper

The Yosys manual, the proc_dff pass; Cummings, "Synthesis and Scripting Techniques for Designing Multi-Asynchronous Clock Designs" and his nonblocking-assignment papers (SNUG).
