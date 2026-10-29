# Next-state and output logic

Once states have codes, the FSM is a register plus two combinational functions: next-state logic from (state bits, inputs) to next state bits, and output logic. Deriving these functions from the state table is mechanical, and seeing how their size depends on the encoding explains why synthesis tools bother to re-encode machines.

## The next-state truth table

With $b$ state bits and $m$ input bits, each next-state bit is a Boolean function of $b + m$ variables, so the next-state truth table has $2^{b+m}$ rows and $b$ output columns. A binary-encoded machine with 6 states and 2 inputs has $b = 3$, so $2^{5} = 32$ rows. Rows for unused state codes are don't-cares unless the design asks for safe recovery.

This table is what the synthesis tool minimizes. Its size grows exponentially in $b + m$, which is one reason large controllers are decomposed into several interacting smaller machines.

## Deriving a next-state bit

Encode four states A = 00, B = 01, C = 10, D = 11 as $s_1 s_0$, with one input $x$. Suppose the transitions on $x = 0, 1$ are A: (B, C), B: (B, D), C: (A, D), D: (C, A). The next value of $s_1$ is 1 when the next state is C or D. Reading row by row over $(s_1, s_0, x)$: A,0 goes to B, $s_1^+ = 0$; A,1 to C, 1; B,0 to B, 0; B,1 to D, 1; C,0 to A, 0; C,1 to D, 1; D,0 to C, 1; D,1 to A, 0.

So $s_1^+$ has four minterms, 001, 011, 101 and 110, and simplifies to $s_1^+ = \bar{s_1}x + \bar{s_0}x + s_1 s_0 \bar{x}$. A different encoding of the same machine gives different functions, and the best one is not obvious in advance.

## One-hot next-state equations

With one-hot encoding the next-state logic can be read straight off the diagram. State bit $X$ is set next cycle exactly when the machine is now in some predecessor state $P$ and the input condition on the arc $P \to X$ holds:

$$X^+ = \sum_{P \to X} P \cdot \text{cond}(P \to X).$$

Each predecessor contributes one term; when both input values lead from $P$ to $X$, the condition is 1 and the term is just $P$. In the machine above, D's predecessors are B (on $x = 1$) and C (on $x = 1$), so $D^+ = B x + C x$, two terms. The logic depth is one AND-OR level regardless of the machine's size, which is why one-hot machines run fast.

## Going deeper

Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), section 3.4; De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 9.
