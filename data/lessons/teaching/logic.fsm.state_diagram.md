# Sequence detectors and state diagrams

Turning an informal specification into a state diagram is the step where control bugs are born. For pattern detectors there is a systematic construction, the same prefix automaton that the Knuth-Morris-Pratt string matcher uses, and it shows how many states a behavior needs and where each transition must go. Formal tools that compile temporal properties into monitor automata perform the same construction.

## How many states

For a detector of a pattern $p$ of length $k$, let state $i$ mean "the longest suffix of the input so far that is a prefix of $p$ has length $i$". A Moore detector needs states $0, 1, \dots, k$, with state $k$ the one that outputs 1: $k + 1$ states. A Mealy detector emits the 1 on the transition that completes the pattern and never needs to sit in state $k$, so it needs $k$ states.

These counts are minimal: from state $i$ the shortest input that produces an output has length $k - i$, a different number for each state, so no two states are equivalent. A detector for 10110 therefore needs 6 Moore states or 5 Mealy states, and needs $\lceil \log_2 6 \rceil = 3$ flip-flops in binary either way.

## Where a mismatch goes

From state $i$, on input $x$: if $x = p_{i+1}$ the machine advances to $i + 1$. Otherwise it falls back to the length of the longest suffix of $p_1 \dots p_i x$ that is also a prefix of $p$, which need not be 0.

For $p = 1011$ in state 3 (matched 101), input 0 gives 1010. Its longest suffix that is a prefix of 1011 is 10, so the next state is 2, not 0. Falling back to 0 instead is the most common bug in hand-written detectors: on the input 101011 the correct machine detects the pattern at the end, while the buggy one, having discarded the 10, misses it. This fallback table is exactly the KMP failure function extended by one input symbol.

## After a detection

In a Moore detector the machine sits in state $k$ after a match. The next input then decides between overlapping and non-overlapping semantics. With overlap, the next state is computed by the same suffix rule from the full pattern: for $p = 1011$ in state 4 with input 0, the string 10110 has longest prefix-suffix 10, so the next state is 2; with input 1, 10111 has longest prefix-suffix 1, so the next state is 1. Without overlap the machine restarts: next state 1 if the input equals $p_1$, else 0.

The number of transitions is the number of states times the number of input values, $(k+1) \cdot 2$ for a one-bit Moore detector. Every one must be specified; a missing transition in RTL becomes an inferred latch or an undefined next state.

## Going deeper

Knuth, Morris and Pratt, "Fast Pattern Matching in Strings" (SIAM Journal on Computing, 1977); Mano and Ciletti, "Digital Design" (5th ed., 2013), chapter 5.
