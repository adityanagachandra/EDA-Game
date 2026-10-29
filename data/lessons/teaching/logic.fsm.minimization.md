# State minimization

Two states are equivalent if no input sequence can tell them apart by the outputs. Merging equivalent states gives a smaller machine with identical behavior. For completely specified machines the minimum is unique and computable in polynomial time by partition refinement, the same algorithm that minimizes deterministic automata; it also underlies sequential equivalence checking, where two machines are equivalent if their product machine has equivalent reset states.

## Equivalence as a fixed point

States $s$ and $t$ of a Moore machine are equivalent when they have the same output and, for every input $x$, their successors $\delta(s, x)$ and $\delta(t, x)$ are equivalent. The definition is recursive, and the coarsest relation satisfying it is computed by refinement:

1. Partition the states by output.
2. Split every block whose members disagree on which block their successors fall in, for some input.
3. Repeat until no block splits.

The final blocks are the equivalence classes, and the minimal machine has one state per class. Each round can only split blocks, so the procedure stops after at most $|S|$ rounds.

## A worked refinement

Take states A to F with outputs A:0, B:0, C:1, D:0, E:1, F:0 and next states on inputs 0 and 1 as A: (B, C), B: (D, E), C: (A, C), D: (B, C), E: (A, E), F: (D, E).

Round 0 blocks by output: $\{A, B, D, F\}$ and $\{C, E\}$. Round 1: for the output-0 block, the successor blocks are A: (0-block, 1-block), B: (0, 1), D: (0, 1), F: (0, 1), all alike, so no split; for $\{C, E\}$: C: (0, 1), E: (0, 1), no split. The partition is stable after one round, with 2 classes. The six-state machine is equivalent to a two-state machine that outputs 1 exactly when the last input was 1.

## Engineering consequences

For a machine with $n$ states and $m$ input values, naive refinement costs $O(n^2 m)$; Hopcroft's algorithm achieves $O(m\, n \log n)$ by always splitting against the smaller half. Synthesis tools apply minimization after extracting an FSM from RTL, but the bigger wins usually come from encoding and logic optimization, since hand-written RTL rarely contains many redundant states.

Incompletely specified machines, with don't-care outputs or transitions, are different: compatibility is not transitive, the minimum is no longer unique, and finding it is NP-hard. In verification, the same refinement idea checks whether two machines are equivalent: build both, refine their disjoint union, and check whether the two reset states end up in the same block.

## Going deeper

Hopcroft, "An n log n Algorithm for Minimizing States in a Finite Automaton" (1971); Kohavi and Jha, "Switching and Finite Automata Theory" (3rd ed., 2010), chapter 10.
