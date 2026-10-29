# FSM state minimization

A finite-state machine written by a designer, or extracted from RTL by a tool, often has redundant states: different states that no input sequence can tell apart. Merging them shrinks the state register and the next-state logic. For completely specified machines the problem is solved exactly and fast by partition refinement, the same algorithm that minimizes deterministic automata.

## Equivalent states and the initial partition

Two states are **equivalent** if, for every input sequence, the machine produces the same output sequence starting from either one. Equivalence is an equivalence relation, so the states split into classes, and the minimal machine has one state per class. Refinement starts from $P_0$: put states in the same block when they produce the same outputs for every single input. In a Moore machine that is grouping by state output; in a Mealy machine by the output pair (output on input 0, output on input 1). For a Mealy machine whose seven states have output pairs (1,1), (1,1), (1,1), (0,1), (1,0), (1,1), (1,1), $P_0$ has 3 blocks.

## Refining until stable

Given partition $P_j$, split each block so that two states stay together only if, for every input, their next states lie in the same block of $P_j$. That gives $P_{j+1}$. Stop when a pass splits nothing: the blocks are then exactly the equivalence classes. Each pass is a sort of the states by the key (current block, block of next state on 0, block of next state on 1). If the first pass splits nothing, the outputs alone already separate every inequivalent pair. A pass can only split blocks, so at most $n - 1$ passes are needed, and Hopcroft's algorithm achieves $O(n \log n)$ by refining with the smaller half of each split.

## Counting passes and building the minimal machine

The number of passes that split something measures how deep the distinctions go: states first separated in $P_j$ need an input sequence of length $j + 1$ to be told apart (for a Mealy machine). If the initial partition is already stable, zero passes split anything. The minimal machine takes one representative per block, with transitions to the blocks of the representatives' successors and the shared outputs. In synthesis, minimization runs before encoding, since fewer states mean fewer flip-flops. For incompletely specified machines, with unspecified outputs or transitions, compatibility is not transitive and exact minimization is NP-hard; tools then use heuristics based on maximal compatibles and covering.

## Going deeper

Moore, "Gedanken-experiments on sequential machines", in Automata Studies (1956); Hopcroft, "An n log n algorithm for minimizing states in a finite automaton" (1971); Kohavi and Jha, "Switching and Finite Automata Theory" (3rd ed., 2010).
