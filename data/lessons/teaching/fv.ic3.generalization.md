# Cube generalization in IC3

The clause IC3 learns when it blocks a state is only as useful as the number of states it excludes. Generalization shrinks the blocked cube literal by literal, keeping each step relatively inductive, so that one clause rules out a large region of the state space. It is where IC3 spends most of its time and where engine implementations differ most, so it is a natural place for a tool designer to tune.

## Why generalize

A blocked state $s$ assigns every latch, so the clause $\neg s$ excludes a single state out of $2^n$. Blocking states one at a time would make IC3 an explicit-state search, hopeless on designs with hundreds of latches. Generalization looks for a sub-cube $c \subseteq s$, with fewer literals, that can be blocked in its place; removing a literal doubles the number of states the clause excludes. A generalized clause must keep the two properties that made $\neg s$ usable: initiation (the cube $c$ contains no initial state) and relative inductiveness ($F_{i-1} \wedge \neg c \wedge T \rightarrow \neg c'$). Generalization never affects soundness, only performance: IC3 with no generalization is still correct, just extremely slow. In practice the reasons a state is unreachable involve a few related latches, so short clauses usually exist.

## Dropping a literal

To test whether a literal can be dropped from cube $s$, form the smaller cube $c$ and run both checks: $c$ must contain no initial state, and no state of the frame $F$ outside $c$ may have a transition into $c$. Example with states $x_2x_1x_0$, initial state 000, $F = \{000, 001, 010\}$, transitions 000 → 001, 001 → 010, 010 → 000 and 010 → 100, and blocked cube $x_2 \wedge x_1 \wedge \neg x_0$ (the state 110). Dropping the $x_2$ literal gives $x_1 \wedge \neg x_0 = \{010, 110\}$; the state 001 of $F$ steps to 010, so it fails. Dropping $x_1$ gives $x_2 \wedge \neg x_0 = \{100, 110\}$; 010 steps to 100, so it fails. Dropping $x_0$ gives $x_2 \wedge x_1 = \{110, 111\}$; no state of $F$ steps into it and it avoids 000, so it passes: the clause $\neg(x_2 \wedge x_1)$ can be learned.

## MIC and unsat cores

Engines generalize in stages. First, the UNSAT proof of the relative-inductiveness query for $\neg s$ often uses only some literals of $s'$; literals whose assumptions are absent from the unsat core can be dropped at once, for free. Then the down or MIC procedure (minimal inductive clause) tries to drop each remaining literal in turn, re-checking with one SAT call per attempt, and keeps a literal only when dropping it breaks the check. Stronger variants handle counterexamples to generalization: when a drop fails because of a predecessor state, ctgDown-style generalization tries to block that predecessor at an earlier frame and then retries the drop. Each refinement buys shorter clauses at the cost of more queries, so the generalization strategy is a policy choice, a trade-off between effort per clause and the number of clauses and frames IC3 needs.

## Going deeper

Bradley, "SAT-Based Model Checking without Unrolling" (VMCAI 2011); Hassan, Bradley and Somenzi, "Better Generalization in IC3" (FMCAD 2013).
