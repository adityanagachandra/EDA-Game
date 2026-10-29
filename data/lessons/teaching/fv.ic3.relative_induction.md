# Relative inductiveness

IC3 builds inductive invariants out of clauses, and the test it applies to each clause is relative inductiveness: a clause need not be inductive on its own, only inductive assuming a frame of already-known facts. This weaker requirement is what lets IC3 learn useful lemmas one at a time. Each check is a single SAT query with one copy of the transition relation, so understanding exactly what it asks is the key to understanding IC3.

## Inductive relative to a frame

Let $F$ be a set of states described by clauses, the frame, and $c$ a clause. $c$ is inductive relative to $F$ if $F \wedge c \wedge T \rightarrow c'$: every transition from a state satisfying both $F$ and $c$ lands in a state satisfying $c$. IC3's clauses usually exclude a cube, a conjunction of literals, so $c = \neg \text{cube}$, and the check fails exactly when some state of $F$ outside the cube steps into the cube. Example with states written $x_2x_1x_0$: $F = \{000, 001, 010\}$, the cube is $\neg x_2 \wedge x_1 \wedge x_0$, the single state 011, and the transitions include 000 → 001, 001 → 011 and 010 → 000. The states of $F \wedge c$ are all of $F$; only 001 steps into the cube, so one state violates consecution and $c$ is not inductive relative to $F$.

## Clauses that exclude larger cubes

A clause with fewer literals excludes a larger cube. The clause $\neg(x_2 \wedge x_0)$ excludes the cube $x_2 \wedge x_0$, which is the two states 101 and 111. The same check applies: count the states of $F$ outside the cube that have a successor inside it. With $F = \{000, 001, 010\}$ and transitions 000 → 001, 001 → 010 and 010 → 101, only 010 steps into the cube, so the count is 1. Larger cubes are more valuable, because one clause rules out more states, but harder to keep relatively inductive, because more states can step into them. Note what is not required: $c$ need not hold on states outside $F$. Assuming $F$ is what makes the query cheap and the lemma easy to find, and it is sound because $F$ over-approximates the states reachable within the current number of steps.

## Initiation and why relative induction is enough

A clause added to IC3's frames must also satisfy initiation: no initial state may lie in the excluded cube, or the clause would be false at reset. With initial state 000, the clause $\neg(\neg x_2 \wedge \neg x_1)$ excludes 000 and 001, so it fails initiation outright, whatever happens in consecution. IC3 checks both conditions for every candidate clause: initiation is a cheap syntactic or SAT check, and consecution is the query $F_{i-1} \wedge c \wedge T \wedge \neg c'$, which must be UNSAT. When both hold, $c$ is added to frames $F_1$ to $F_i$. Since $F_{i-1}$ contains every state reachable within $i - 1$ steps, relative inductiveness guarantees that $c$ holds on every state reachable within $i$ steps. When the frames stop changing, all their clauses are mutually inductive, which is the final proof.

## Going deeper

Bradley, "SAT-Based Model Checking without Unrolling" (VMCAI 2011); Een, Mishchenko and Brayton, "Efficient Implementation of Property Directed Reachability" (FMCAD 2011).
