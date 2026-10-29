# Distinguishing sequences

Two inequivalent states are told apart by some input sequence that makes their outputs differ. The shortest such sequence is the evidence a tool reports when it refuses to merge states, and it is what a sequential equivalence checker returns as a counterexample. Its length is tied directly to the passes of partition refinement.

## One-step distinguishable pairs

In a Mealy machine, states $p$ and $q$ are distinguished by a single input if their outputs differ on input 0 or on input 1. These are exactly the pairs in different blocks of $P_0$. Count them by comparing the output pairs of every two states. If five states have output pairs (1,1), (0,0), (0,0), (1,1), (0,0), the pairs split by $P_0$ are those mixing a (1,1) state with a (0,0) state, $2 \cdot 3 = 6$ pairs. Pairs inside a block of $P_0$ need longer sequences, or are equivalent.

## Breadth-first search over state pairs

To find the shortest distinguishing sequence for $p$ and $q$, search over pairs of states. Start at $(p, q)$. From a pair $(s, t)$, an input $x$ either produces different outputs, ending the search at the current depth plus one, or leads to the pair $(\delta(s, x), \delta(t, x))$. Breadth-first search visits pairs in order of sequence length, so the first difference found gives the shortest length. If the search exhausts the reachable pairs without a difference, the states are equivalent. The search space is at most $n^2$ pairs, so this is polynomial, and recording the inputs along the way gives the sequence itself.

## Length bounds and refinement

A pair first separated in partition $P_j$ needs a sequence of length exactly $j + 1$ in a Mealy machine: the pair stays together in $P_0, \dots, P_{j-1}$ precisely because no shorter sequence separates it. So the longest shortest-distinguishing sequence over all inequivalent pairs equals the number of refinement steps counting $P_0$ as step 1. Since a machine with $n$ states has at most $n$ blocks and each step adds one, any two inequivalent states are separated by a sequence of length at most $n - 1$. Sequential equivalence checkers and FSM test generators rely on this bound, and on the same pair-graph search.

## Going deeper

Kohavi and Jha, "Switching and Finite Automata Theory" (3rd ed., 2010), on state equivalence and distinguishing sequences; Lee and Yannakakis, "Principles and methods of testing finite state machines: a survey", Proceedings of the IEEE (1996).
