# Reset, unreachable and illegal states

An FSM's specification describes its legal states; its hardware has $2^b$ codes. Reset puts the machine into a known legal state, reachability analysis finds which states can actually occur, and the remaining codes are illegal. How a design and its tools treat these three sets decides robustness to soft errors, the size of the logic, and what a formal tool must prove.

## Reset establishes the initial state

At power-up, flip-flops hold arbitrary values. Reset forces the state register to the designated initial code. A synchronous reset is part of the next-state logic and acts at a clock edge; an asynchronous reset acts immediately through dedicated flip-flop pins and needs its release synchronized. Formal verification and sequential equivalence checking start from the reset state, so a missing or incomplete reset (some flip-flops not reset) makes the initial state set larger, and properties that rely on initialization may fail.

Not every flip-flop needs reset: datapath registers that are always written before being read can skip it, saving area and routing. Control state, however, must always be reset.

## Reachable versus unreachable states

A state is reachable if some input sequence leads to it from reset. Reachability is a graph search from the reset state over the transition relation, cheap for an explicit machine and the core of model checking for a symbolic one, where the state set is represented as a BDD or explored with SAT.

Unreachable states can arise from the specification (a state no transition enters) or from the encoding (unused codes). Synthesis may exploit unreachability: if a code never occurs, its next-state rows are don't-cares, which can shrink logic substantially. Sequential optimizations that rely on unreachable states, however, produce designs that are equivalent only from reset, so the equivalence checker must be told the reset state to prove them correct.

## Illegal states and recovery

An illegal state is a code outside the legal set, reachable only through a fault such as a single-event upset or a timing violation. Two policies exist. Treating illegal codes as don't-cares gives the smallest logic, but the machine may then lock up in a cycle of illegal states forever. A safe FSM adds logic so every illegal code transitions to a known state (often reset or an error state), and may add detection that raises an alarm.

The choice is a reliability decision, common in automotive, aerospace and safety-critical designs, and synthesis tools expose it as an FSM attribute. A formal tool can check the safe policy directly: from any code, reset or recovery is reached within a bounded number of cycles.

## Going deeper

Kohavi and Jha, "Switching and Finite Automata Theory" (3rd ed., 2010); Cummings, "The Fundamentals of Efficient Synthesizable Finite State Machine Design" (SNUG 2003).
