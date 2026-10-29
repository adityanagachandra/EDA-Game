# Bounded model checking

Bounded model checking (BMC) asks whether a bad state is reachable within $k$ steps by unrolling the circuit into $k + 1$ copies, one per clock cycle, and handing the result to a SAT solver. It is the most widely used formal engine for finding bugs, because it needs no BDDs, scales to large designs and returns a shortest counterexample. Estimating the size of the unrolled formula is the first step in predicting whether a BMC run will finish.

## Unrolling a circuit into frames

Take a circuit with $n$ latches, $m$ primary inputs and $g$ two-input AND gates, with one gate output as the bad signal. Frame $i$ is a copy of the combinational logic for cycle $i$, for $i = 0, \dots, k$. The convention here: frame 0's latch outputs are fresh variables, constrained to the reset values; in every later frame each latch output is simply the previous frame's next-state gate output, so it needs no variable of its own; and every frame has its own copy of the inputs and the gates. The variable count is therefore $n + (k + 1)(m + g)$. With $n = 2$, $m = 1$, $g = 4$ and $k = 3$: $2 + 4 \cdot (1 + 4) = 22$ variables. Real tools such as ABC's bmc engines unroll AIGs in exactly this frame-by-frame way.

## Counting the clauses

The CNF has three parts. Initial-state constraints: one unit clause per latch fixing its reset value, $n$ clauses. Transition logic: each 2-input AND gate costs 3 Tseitin clauses per frame, so $3g(k + 1)$. Property: one clause $(\text{bad}_0 \vee \text{bad}_1 \vee \dots \vee \text{bad}_k)$ saying the bad signal is 1 in some frame. Total $n + 3g(k + 1) + 1$. In the example, $2 + 3 \cdot 4 \cdot 4 + 1 = 51$ clauses. A satisfying assignment is a counterexample: read the input variables of frames 0 to $k$ as a stimulus and the latch values as the state trace. If the bad signal is 1 in frame $j$, the trace has length $j$, and replaying it in simulation confirms the failure independently of the solver.

## Growing the bound

BMC usually runs for $k = 0, 1, 2, \dots$ until it finds a counterexample or reaches a time or depth limit. Going from depth $k$ to $k + 1$ adds one frame: $m + g$ variables and $3g$ gate clauses. The property clause gains one literal but remains one clause, and the reset clauses do not change. For $g = 4$ that is 12 clauses per extra frame; for a million-gate design, 3 million. The formula grows linearly with depth, but the search can grow much faster, which is why deep bounds are expensive. The crucial limitation: an UNSAT answer at depth $k$ only proves that no counterexample of length up to $k$ exists. It says nothing about longer runs, so BMC alone is a bug-finding engine; proofs need a completeness threshold, induction or IC3.

## Going deeper

Biere, Cimatti, Clarke and Zhu, "Symbolic Model Checking without BDDs" (TACAS 1999); Biere, Cimatti, Clarke, Strichman and Zhu, "Bounded Model Checking" (Advances in Computers, 2003).
